#!/usr/bin/python

import label
import common as c
import matplotlib.pyplot as plt
import numpy as np

def plot_bar(output_file, y_label, y_l, y_err, x_label, x_l):
  fs = 14; bw = 0.5; lw = 1.5; cs=2.5; ms=4
  plt.figure(figsize=(8,4))
  idx = np.arange(0, len(x_l))
  plt.bar(idx, y_l, bw, yerr=y_err, edgecolor='black', capsize=cs, align='center')
  plt.ylabel(y_label, fontsize=fs)
  plt.ylim(ymin=0)
  plt.xlabel(x_label, fontsize=fs)
  plt.xticks(idx, x_l, fontsize=fs)
  plt.tick_params(axis='both',direction='in',labelsize=fs)
  plt.savefig(output_file, bbox_inches='tight')
