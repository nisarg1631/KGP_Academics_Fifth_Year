# COMPUTATIONAL BIOPHYSICS PROJECT
## Project Title: Identification of  Disease-Associated Genes using ML

### Requirements
1. Scikit-learn
2. Pandas
3. Numpy
4. Matplotlib
5. Tensorflow

### Files
1. disgenet_2020.db: Database file of DisGeNET
2. gene_sequences_final.{csv, fasta}: Gene sequences to be used in training downloaded from NCBI
3. model.py: Contains the code for training the model using traditional ML algorithms
4. model_neural.py: Contains the code for training the model using neural networks
5. features/ : directory containing the feature vectors for the seuqences in gene_sequences_final.csv
6. model_pca_10/ : directory containing the trained models using PCA with 10 components
7. model_pca_20/ : directory containing the trained models using PCA with 20 components
8. model_neural/ : directory containing the trained neural network model
9. results/ : directory containing the results of the models

### Note
The files disgenet_2020.db, gene_sequences_final.{csv, fasta}, features/ and models/ directories have not been included in the submission due to their large size. They can be downloaded from the GIT repository at https://github.com/nisarg1631/CBP-Project
