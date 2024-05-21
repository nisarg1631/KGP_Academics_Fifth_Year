import sqlite3
import pandas as pd
import pickle

conn = sqlite3.connect('disgenet_2020.db')
cursor = conn.cursor()

cursor.execute("SELECT geneNID, diseaseNID FROM geneDiseaseNetwork GROUP BY geneNID, diseaseNID;")
gene2disease = cursor.fetchall()

cursor.execute("SELECT diseaseNID, diseaseClassNID FROM disease2class;")
disease2class = { int(c[0]):int(c[1]) for c in cursor.fetchall() }

# create one hot encoding for each gene to the disease classes it belongs to
gene2class = {}
for gene, disease in gene2disease:
    if disease in disease2class:
        if gene not in gene2class:
            gene2class[gene] = [disease2class[disease]]
        elif disease2class[disease] not in gene2class[gene]:
            gene2class[gene].append(disease2class[disease])

# map the disease classes to a unique index
class2index = { c:i for i, c in enumerate(set([c for classes in gene2class.values() for c in classes])) }
print(class2index)

# get the class names from the diseaseClass table
cursor.execute("SELECT diseaseClassNID, diseaseClass, diseaseClassName FROM diseaseClass;")
class2name = { int(c[0]):(c[1], c[2]) for c in cursor.fetchall() }

print({class2index[c]:class2name[c] for c in class2index})

# create one hot encoding for each gene
gene2class_onehot = {}
for gene, classes in gene2class.items():
    onehot = [0] * len(class2index)
    for c in classes:
        onehot[class2index[c]] = 1
    gene2class_onehot[gene] = onehot

# read the sequences from gene_sequences_final.csv no header
df_seq = pd.read_csv('gene_sequences_final.csv', header=None)

# read the features for each sequence from the features directory - ent.csv, oth_cleaned.csv, aad_cleaned.csv, acd_cleaned.csv, soc_cleaned.csv
# add the features to the corresponding row in df_seq

df_ent = pd.read_csv('features/ent.csv', header=None)
df_oth = pd.read_csv('features/oth_cleaned.csv', header=None)
df_aad = pd.read_csv('features/aad_cleaned.csv', header=None)
df_acd = pd.read_csv('features/acd_cleaned.csv', header=None)
df_soc = pd.read_csv('features/soc_cleaned.csv', header=None)

# merge all the dataframes row wise
df_seq = pd.concat([df_seq, df_ent, df_oth, df_aad, df_acd, df_soc], axis=1)

# name columns sequentially from 0 to len
df_seq.columns = range(len(df_seq.columns))

# add the one hot encoding for the genes which are present for others drop them
df_classes = df_seq.iloc[:, 0].map(lambda x: gene2class_onehot.get(x, [0]*len(class2index)))
df_seq['classes'] = df_classes
df_seq = df_seq[df_seq['classes'].map(lambda x: sum(x) > 0)]

# df_classes = df_seq.iloc[:, 0].map(lambda x: gene2class_onehot.get(x, [0]*len(class2index)))
# df_seq['classes'] = df_classes

print(df_seq.head())

# drop the first two columns and shuffle the rows
df_seq = df_seq.drop(columns=[0, 1])
df_seq = df_seq.sample(frac=1).reset_index(drop=True)

print(df_seq.head())

# build a neural network model
# input layer: all feature columns (1094)
# output layer: one hot encoding of the classes (26)
# hidden layers: select appropriate number of layers and neurons
# activation function: relu for hidden layers, softmax for output layer
# loss function: categorical crossentropy
# optimizer: adam
# metrics: accuracy for each class

import tensorflow as tf
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import Dense, Input
from tensorflow.keras.optimizers import Adam
from tensorflow.keras.losses import BinaryCrossentropy
from sklearn.model_selection import train_test_split
from sklearn.metrics import multilabel_confusion_matrix, classification_report

X = df_seq.iloc[:, :-1]
y = df_seq.iloc[:, -1]

# train test val split 60-20-20
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.4)
X_test, X_val, y_test, y_val = train_test_split(X_test, y_test, test_size=0.5)

X_train = tf.convert_to_tensor(X_train.values, dtype=tf.float32)
X_val = tf.convert_to_tensor(X_val.values, dtype=tf.float32)

y_train = tf.convert_to_tensor(y_train.values.tolist(), dtype=tf.float32)
y_val = tf.convert_to_tensor(y_val.values.tolist(), dtype=tf.float32)

model = Sequential([
    Input(shape=(1094,)),
    Dense(512, activation='relu'),
    Dense(256, activation='relu'),
    Dense(128, activation='relu'),
    Dense(26, activation='sigmoid')
])


# model.compile(optimizer=Adam(), loss=BinaryCrossentropy())
# history = model.fit(X_train, y_train, epochs=50, batch_size=32, validation_data=(X_val, y_val))
# history = history.history

# model.save('model_neural_2.keras')
# with open('model_neural_2_history.pkl', 'wb') as f:
#     pickle.dump(history, f)

# load model
model = tf.keras.models.load_model('model_neural_2.keras')

# load history
with open('model_neural_2_history.pkl', 'rb') as f:
    history = pickle.load(f)

# display the epoch vs loss and val_loss curve with fonts size 22
import matplotlib.pyplot as plt

plt.plot(history['loss'])
plt.plot(history['val_loss'])
plt.title('Model loss', fontsize=22)
plt.ylabel('Loss', fontsize=22)
plt.xlabel('Epoch', fontsize=22)
plt.legend(['Train', 'Validation'], loc='upper left', fontsize=22)

# increase ticks and legend font to 22
plt.xticks(fontsize=22)
plt.yticks(fontsize=22)

plt.show()

y_pred = model.predict(X_test)
print(y_pred)

# get the labels from the probabilities with threshold 0.1,0.2, 0.5
labels_01 = tf.cast(y_pred > 0.1, tf.bool)
labels_02 = tf.cast(y_pred > 0.2, tf.bool)
labels_05 = tf.cast(y_pred > 0.5, tf.bool)

print(labels_01)
y_test = y_test.values.tolist()
print(y_test)

# get the confusion matrix and classification report
conf_matrix_01 = multilabel_confusion_matrix(y_test, labels_01)
conf_matrix_02 = multilabel_confusion_matrix(y_test, labels_02)
conf_matrix_05 = multilabel_confusion_matrix(y_test, labels_05)

report_01 = classification_report(y_test, labels_01)
report_02 = classification_report(y_test, labels_02)
report_05 = classification_report(y_test, labels_05)

print("Confusion matrix 0.1: ", conf_matrix_01)
print("Classification report 0.1: ", report_01)

print("Confusion matrix 0.2: ", conf_matrix_02)
print("Classification report 0.2: ", report_02)

print("Confusion matrix 0.5: ", conf_matrix_05)
print("Classification report 0.5: ", report_05)
