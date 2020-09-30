#!/usr/bin/python
import matplotlib
import common as c

#matplotlib.rcParams['ps.useafm'] = True
#matplotlib.rcParams['pdf.use14corefonts'] = True
#matplotlib.rcParams['text.usetex'] = True
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42

stat_label = ['Mean', 'Median', '95%', '99%']
metric_label = {
    c.Stat.MEAN : 'Mean', c.Stat.MEDIAN : 'Median', \
    c.Stat.P95th:'95th%', c.Stat.P99th:'99th%', \
    }
protocol_name = {
    'do':'do', 'h':'hybrid', \
    'd':'dynamic', 'dp':'dynamic-paxos', 'dfp':'dynamic-fp', \
    'dt':'dynamic-to', 'dfpt':'dynamic-fp-to', \
    'm': 'mencius', 'mc': 'mencius-commit', 'e':'epaxos', 'ef':'epaxos-thrifty', \
    'p':'paxos', 'pf':'paxos-thrifty', \
    'fp':'fastpaxos',\
    }
protocol_label = {
    'do':'DO-FastPaxos', 'h':'Hybrid', \
    'd':'HC', 'dp':'HC-VMencius', 'dfp':'HC-FastPaxos', \
    'dt':'Domino', 'dfpt':'HC-FastPaxos-TO', 'dtsync':'HC-TO-sync', \
    'm': 'Mencius', 'mc': 'Mencius-Commit', \
    'e': 'EPaxosNoThrifty', 'ef':'EPaxos', \
    'p':'Multi-Paxos', 'pf':'Multi-PaxosThrifty', \
    'fp': 'Fast Paxos', \
    }
line_style = {
    'do' : '--', 'h' : '--', \
    'd'  : '-.', 'dp':'--', 'dfp':'-', \
    'dt' : '-', 'dfpt':'-.', 'dtsync':'--', \
    'm'  : '--', 'mc': ':', \
    'ef' : ':',  'e' : '-.', \
    'p'  : '-.', 'pf': '--', \
    'fp' : '-', \
    }

line_color = {
    'do' : 'darkgreen', 'h' : 'gold', \
    'd'  : 'darkgoldenrod', 'dp' : 'gold', 'dfp' : 'green', \
    'dt' : 'black', 'dfpt':'black', 'dtsync':'black', \
    'm'  : 'darkgreen', 'mc': 'forestgreen', \
    'ef' : 'firebrick', 'e' : 'red', \
    'p'  : 'darkblue', 'pf' : 'blue', \
    'fp' : 'grey', \
    }
point_fmt = {
    'do' : '<', 'h' : '>', \
    'd': '^', 'dp':'<', 'dfp':'>', \
    'dt': 'v', 'dfpt':'2', 'dtsync':'3',\
    'm' : 's', 'mc' : 'x', \
    'ef': 'o', 'e':'P', \
    'p' : '*', 'pf' : '1' , \
    'fp' : '>', \
    }
point_color = {
    'do' : None, 'h' : None, \
    'd' : None, 'dp' : None, 'dfp' : None, \
    'dt': None, 'dfpt' : None, 'dtsync':None, \
    'm' : None, 'mc' : None, \
    'ef' : 'white', 'e' : None, \
    'p' : None, 'pf' : None,\
    'fp' : None, \
    }
err_color = 'darkblue'
#Hatch Type: '', '-', '+', 'x', '\\', '*', 'o', 'O', '.', '/'
hatch_type = {
    'do' : '.', 'h' : 'o', \
    'd' : '', 'dp' : '\\', 'dfp' : '/', \
    'dt': '', 'dfpt':'O', \
    'm' : 'o', 'mc' : '.', \
    'ef': '/', 'e' : '\\', \
    'p' : 'x', 'pf' : '+', \
    'fp': '/', \
    }
bar_color = {
    'do' : 'white', 'h' : 'white', \
    'd' : 'black', 'dp' : 'grey', 'dfp' : 'grey', \
    'dt': 'white', 'dfpt':'grey', \
    'm' : 'white', 'mc' : 'grey', \
    'ef': 'white', 'e'  : 'grey', \
    'p' : 'white', 'pf' : 'grey', \
    'fp': 'white', \
    }
