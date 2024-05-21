import sqlite3
import pandas as pd

conn = sqlite3.connect('disgenet_2020.db')
cursor = conn.cursor()

# in the table geneDiseaseNetwork aggregate based on geneNID and average the score
cursor.execute("SELECT geneNID, SUM(score) FROM geneDiseaseNetwork GROUP BY geneNID;")
score_dict = { int(c[0]):float(c[1]) for c in cursor.fetchall() }

# read the sequences from gene_sequences_final.csv no header
df_seq = pd.read_csv('gene_sequences_final.csv', header=None)

# read the features for each sequence from the features directory - ent.csv, oth_reduced.csv, aad_reduced.csv, acd_reduced.csv, soc_reduced.csv
# add the features to the corresponding row in df_seq

df_ent = pd.read_csv('features/ent.csv', header=None)
df_oth = pd.read_csv('features/oth_reduced.csv', header=None)
df_aad = pd.read_csv('features/aad_reduced.csv', header=None)
df_acd = pd.read_csv('features/acd_reduced.csv', header=None)
df_soc = pd.read_csv('features/soc_reduced.csv', header=None)

# merge all the dataframes row wise
df_seq = pd.concat([df_seq, df_ent, df_oth, df_aad, df_acd, df_soc], axis=1)

# name columns sequentially from 0 to len
df_seq.columns = range(len(df_seq.columns))

df_scores = df_seq.iloc[:, 0].map(lambda x: score_dict.get(x, 0))
df_isdisease = df_seq.iloc[:, 0].map(lambda x: 1 if x in score_dict else 0)

# counts
print(df_isdisease.value_counts())

# df_seq['score'] = df_scores
df_seq['is_disease'] = df_isdisease

print(df_seq.head())

# drop the first two columns and shuffle the rows
df_seq = df_seq.drop(columns=[0, 1])
df_seq = df_seq.sample(frac=1).reset_index(drop=True)

print(df_seq.head())

# split the data into training and testing data
from sklearn.model_selection import train_test_split
from sklearn.svm import SVC
from sklearn.ensemble import GradientBoostingClassifier
from sklearn.ensemble import RandomForestClassifier
from sklearn.ensemble import AdaBoostClassifier
from sklearn.discriminant_analysis import QuadraticDiscriminantAnalysis
from sklearn.metrics import accuracy_score
from sklearn.metrics import confusion_matrix
from sklearn.metrics import classification_report
import joblib

X = df_seq.iloc[:, :-1]
y = df_seq.iloc[:, -1]

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.25)

# train SVM and save it to a file
# model = SVC(verbose=True)
# model.fit(X_train, y_train)
# joblib.dump(model, 'modelSVC.pkl')

# load the model from the file
model = joblib.load('modelSVC.pkl')

# test the model and print accuracy, f1-score and confusion matrix
y_pred = model.predict(X_test)
accuracy = accuracy_score(y_test, y_pred)
conf_matrix = confusion_matrix(y_test, y_pred)
report = classification_report(y_test, y_pred)

print("SVM model")
print("Accuracy: ", accuracy)
print("Confusion matrix: ", conf_matrix)
print("Classification report: ", report)

# train GradientBoostingClassifier and save it to a file
# model = GradientBoostingClassifier(verbose=True)
# model.fit(X_train, y_train)
# joblib.dump(model, 'modelGBC.pkl')

# load the model from the file
model = joblib.load('modelGBC.pkl')

# test the model and print accuracy, f1-score and confusion matrix
y_pred = model.predict(X_test)
accuracy = accuracy_score(y_test, y_pred)
conf_matrix = confusion_matrix(y_test, y_pred)
report = classification_report(y_test, y_pred)

print("GradientBoostingClassifier model")
print("Accuracy: ", accuracy)
print("Confusion matrix: ", conf_matrix)
print("Classification report: ", report)

# train Random Forest Classifier and save it to a file
# model = RandomForestClassifier(verbose=True, n_jobs=-1)
# model.fit(X_train, y_train)
# joblib.dump(model, 'modelRFC.pkl')

# load the model from the file
model = joblib.load('modelRFC.pkl')

# test the model and print accuracy, f1-score and confusion matrix
y_pred = model.predict(X_test)
accuracy = accuracy_score(y_test, y_pred)
conf_matrix = confusion_matrix(y_test, y_pred)
report = classification_report(y_test, y_pred)

print("RandomForestClassifier model")
print("Accuracy: ", accuracy)
print("Confusion matrix: ", conf_matrix)
print("Classification report: ", report)

# train ADABoost Classifier and save it to a file
# model = AdaBoostClassifier()
# model.fit(X_train, y_train)
# joblib.dump(model, 'modelABC.pkl')

# load the model from the file
model = joblib.load('modelABC.pkl')

# test the model and print accuracy, f1-score and confusion matrix
y_pred = model.predict(X_test)
accuracy = accuracy_score(y_test, y_pred)
conf_matrix = confusion_matrix(y_test, y_pred)
report = classification_report(y_test, y_pred)

print("AdaBoostClassifier model")
print("Accuracy: ", accuracy)
print("Confusion matrix: ", conf_matrix)
print("Classification report: ", report)

# train Quadratic Discriminant Analysis and save it to a file
# model = QuadraticDiscriminantAnalysis()
# model.fit(X_train, y_train)
# joblib.dump(model, 'modelQDA.pkl')

# load the model from the file
model = joblib.load('modelQDA.pkl')

# test the model and print accuracy, f1-score and confusion matrix
y_pred = model.predict(X_test)
accuracy = accuracy_score(y_test, y_pred)
conf_matrix = confusion_matrix(y_test, y_pred)
report = classification_report(y_test, y_pred)

print("QuadraticDiscriminantAnalysis model")
print("Accuracy: ", accuracy)
print("Confusion matrix: ", conf_matrix)
print("Classification report: ", report)
