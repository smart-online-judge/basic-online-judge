import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd
import numpy as np

def heatmap(x_labels, y_labels, values):
    fig, axes = plt.subplots()
    _ = axes.imshow(values)

    axes.set_xticks(np.arange(len(x_labels)))
    axes.set_yticks(np.arange(len(y_labels)))

    axes.set_xticklabels(x_labels)
    axes.set_yticklabels(y_labels)

    # Rotate the tick labels and set their alignment.
    plt.setp(axes.get_xticklabels(), rotation=45, ha="right", fontsize=10,
             rotation_mode="anchor")

    for i in range(len(y_labels)):
        for j in range(len(x_labels)):
            axes.text(j, i, "%.2f"%values[i, j],
                      ha="center", va="center", color="w", fontsize=6)

    fig.tight_layout()
    plt.show()

def heatmap2(x_labels, y_labels, values):
    df = pd.DataFrame(values)
    df.columns = x_labels
    df.index = y_labels

    fig, ax = plt.subplots(figsize=(5,5))
    # Rotate the tick labels and set their alignment.
    plt.setp(ax.get_xticklabels(), rotation=45, ha="right", fontsize=10,
             rotation_mode="anchor")
    ax_blue = sns.heatmap(df, cmap="YlGnBu")
    # ax_red = sns.heatmap(df)
    plt.yticks(rotation = 0)
    plt.show()