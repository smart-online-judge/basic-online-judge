import numpy as np

def get_averaged_similarity_mat(syn_mat, sem_mat):
    syn_weights = np.full(syn_mat.shape, 0.6)
    sem_weights = np.full(sem_mat.shape, 0.4)

    average_mat = np.average(
        np.array([syn_mat, sem_mat]),
        axis=0,
        weights=np.array([syn_weights, sem_weights]))

    average_mat = np.round(average_mat, 3)
    return average_mat.tolist()