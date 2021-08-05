#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import os
import re
import math
import numpy as np
import json

def get_gaussian_td(fname) :
#    f = open('%s/%s.td' % (fname, fname), 'r')
    f = open('%s' % (fname), 'r')
    E0 = 0.0
    H = 0.0
    S = 0.0
    G = 0.0
    for line in f :
	s = line.rstrip('\r\n')
	if 'Sum of electronic and zero-point Energies=' in s :
	    i = s.index("=")+1
	    E0 = float(s[i:])
	if 'Sum of electronic and thermal Enthalpies=' in s :
	    i = s.index("=")+1
	    H = float(s[i:])
	if 'Sum of electronic and thermal Free Energies=' in s :
	    i = s.index("=")+1
	    G = float(s[i:])
	    S = G - H
###	    print 'E0=%.8f  H=%.8f  S=%.8f  G=%.8f ' % (E0, H, S, G)
#	    break
    return (E0, H, S, G)


def get_gaussian_wiberg(fname) :
    mr = []
    mfreq = []
    f = open('%s' % (fname), 'r')
    line = f.readline()
    flag = False
    NAtoms = 0
    while line:
#	print line
	line = f.readline()
	s = line.rstrip('\r\n')
	if ' NAtoms=' in s : 
		NAtoms = int(re.split(r'\s+', s)[2])
#	if 'E(B3LYP)' in s : print s[1:]
#	if 'E(UB3LYP)' in s : print s[1:]
#	if 'Zero-point correction=' in s : print s[1:]
#	if 'Thermal correction to Energy=' in s : print s[1:]
#	if 'Thermal correction to Enthalpy=' in s : print s[1:]
#	if 'Thermal correction to Gibbs Free Energy=' in s : print s[1:]
	if 'Sum of electronic and zero-point Energies=' in s : print s[1:]
	if 'Sum of electronic and thermal Energies=' in s : print s[1:]
	if 'Sum of electronic and thermal Enthalpies=' in s : print s[1:]
	if 'Sum of electronic and thermal Free Energies=' in s : print s[1:]
#	print s
	if ' Wiberg bond index matrix in the NAO basis:' in s :
	    f.readline()
	    f.readline()
	    f.readline()
	    flag = True
	    mw = []
	    for i in range(NAtoms) : mw.append([])
#	    print ' *** init ', mw
	    continue
	if 'Wiberg bond index, Totals by atom:' in s :
	    flag = False
#	    print ' *** Wiberg\n', mw
	    mr.append(mw)
	    continue
	if flag and len(s)<2 : continue
	if flag and 'Atom' in s : continue
	if flag and '----' in s : continue
	if flag :
	    i = int(re.split(r'\s+', s)[1][:-1])
	    m = list(map(float, re.split(r'\s+', s)[3:]))
#	    print ' ***', i, mw[i-1], type(mw[i-1])
	    mw[i-1].extend(m)
#	    print ' ***', i, mw[i-1]
	if ' Frequencies --' in s :
	    freq = list(map(float, re.split(r'\s+', s[16:])[1:]))
#	    print 'freq: ', freq, re.split(r'\s+', s[16:])
	    mfreq.extend(freq)
	    continue
    f.close()
    return mr[-3], mfreq


# Wiberg bond index matrix in the NAO basis:                                    
#
#     Atom    1       2       3       4       5
#     ---- ------  ------  ------  ------  ------
#   1.  C  0.0000  1.7840  1.0450  0.0050  0.9055
#   2.  O  1.7840  0.0000  0.1664  0.0096  0.0535
#   3.  O  1.0450  0.1664  0.0000  0.7179  0.0123
#   4.  H  0.0050  0.0096  0.7179  0.0000  0.0177
#   5.  H  0.9055  0.0535  0.0123  0.0177  0.0000

# Frequencies --    622.5210               677.2119              1045.6598

#E(B3LYP)= -1235.94645964
#Zero-point correction= 0.229019
#Thermal correction to Energy= 0.247022
#Thermal correction to Enthalpy= 0.247966
#Thermal correction to Gibbs Free Energy= 0.177937
#Sum of electronic and zero-point Energies= -1235.717441
#Sum of electronic and thermal Energies= -1235.699438
#Sum of electronic and thermal Enthalpies= -1235.698494
#Sum of electronic and thermal Free Energies= -1235.768523

# NAtoms=     19

# Mulliken charges and spin densities:
#               1          2
#     1  C    0.242763   0.000000
#    31  H   -0.430101  -0.000000
# Sum of Mulliken charges =  -0.00000  -0.00000

# Mulliken charges:
#               1
#     1  C    0.648587
#    35  H   -0.377221
# Sum of Mulliken charges =  -0.00000

def read_charge(name, chtype):
    if chtype == 'Mulliken' : flagstr = ' %s charges and spin densities:' % chtype
    else : flagstr = ' %s charges:' % chtype
    flagend = ' Sum of %s' % chtype
    n = 0
    charges = []
    atoms = []
    m = []
    with open( '%s' % (name), 'r') as f:
	line = f.readline()
	while line:
	    if line.startswith(flagstr) :
#		print '[%s]' % line
		n = 0
		charges = []
		atoms = []
		m = []
		line = f.readline()
		line = f.readline()
		while not flagend in line :
		    lst = line.strip().split()
#		    print lst
#		    if len(lst)>1 : m[].append(float(lst)
		    if len(lst)>1 :
			charges.append(float(lst[2]))
			atoms.append(lst[1])
			m.append(lst[0])
			n = n+1
		    line = f.readline()
	    line = f.readline()
    f.close()
#    print charges
    return charges, atoms, m

def read_charge_NBO(name):
    flagstr = ' Summary of Natural Population Analysis:'
    flagend = ' ================'
    n = 0
    charges = []
    oldcharges = list(charges)
    oldoldcharges = list(oldcharges)
    atoms = []
    m = []
    with open( '%s' % (name), 'r') as f:
	line = f.readline()
	while line:
	    if line.startswith(flagstr) :
		line = f.readline()
		line = f.readline()
		line = f.readline()
		line = f.readline()
		line = f.readline()
		n = 0
		oldoldcharges = list(oldcharges)
		oldcharges = list(charges)
		charges = []
		atoms = []
		m = []
		while not flagend in line :
		    lst = line.strip().split()
#		    print lst
#		    if len(lst)>1 : m[].append(float(lst)
		    if len(lst)>1 : 
			charges.append(float(lst[2]))
			atoms.append(lst[1])
			m.append(lst[0])
			n = n+1
		    line = f.readline()
	    line = f.readline()
    f.close()
#    print(charges)
#    print(oldcharges)
#    print(oldoldcharges)
    return oldoldcharges, atoms, m


def read_opt_geom(name):
    flag1 = 0
    flagstr1 = ' Optimization completed.'
    flagstr2 = '                          Input orientation:'
    flagstr3 = '                         Standard orientation:'
    flagend = ' ---------------------------------------------------------------------'
    with open( '%s' % (name), 'r') as f:
	line = f.readline()
	while line:
	    if line.startswith(flagstr1) : flag1=1
	    if flag1==1 and line.startswith(flagstr2) :
		line = f.readline()
		line = f.readline()
		line = f.readline()
		line = f.readline()
		line = f.readline()
		atoms = []
		while not flagend in line :
		    lst = line.strip().split()
		    print line
		    if len(lst)>1 : atoms.append(lst)
		    line = f.readline()
	    if flag1==1 and line.startswith(flagstr3) :
		line = f.readline()
		line = f.readline()
		line = f.readline()
		line = f.readline()
		line = f.readline()
		satoms = []
		while not flagend in line :
		    lst = line.strip().split()
		    print line
		    if len(lst)>1 : satoms.append(lst)
		    line = f.readline()
	    line = f.readline()
    f.close()
    print atoms
    print satoms
    return atoms, satoms

# Optimization completed.
#    -- Stationary point found.
###
#                          Input orientation:                          
# ---------------------------------------------------------------------
# Center     Atomic      Atomic             Coordinates (Angstroms)
# Number     Number       Type             X           Y           Z
# ---------------------------------------------------------------------
#      1          7           0       -2.598568   -0.531148   -0.105067
#      2          7           0       -3.701252   -0.574992   -0.140603
# ---------------------------------------------------------------------
#                         Standard orientation:                         
# ---------------------------------------------------------------------
# Center     Atomic      Atomic             Coordinates (Angstroms)
# Number     Number       Type             X           Y           Z
# ---------------------------------------------------------------------
#      1          7           0        0.000000    0.000000    0.552063
#      2          7           0        0.000000   -0.000000   -0.552063
# ---------------------------------------------------------------------
# Rotational constants (GHZ):           0.0000000          59.2087908          59.2087908






# main

sysname = sys.argv[1]
outname = sys.argv[2]

E0, H, S, G = get_gaussian_td(sysname)
mw, mfreq = get_gaussian_wiberg(sysname)

#print mw
#print mfreq

###print '\nWiberg:'
###for lst in mw :
###    for w in lst:
###	print '%8.4f' % (w),
###    print ''

###print '\nFreq:'
###for freq in mfreq :
###    print '%10.4f' % (freq),
###print ''

Ch_NBO, Atoms, M = read_charge_NBO(sysname)
Ch_APT, Atoms, M = read_charge(sysname, 'APT')
Ch_ESP, Atoms, M = read_charge(sysname, 'ESP')
Ch_Mulliken, Atoms, M = read_charge(sysname, 'Mulliken')
Geom, SGeom = read_opt_geom(sysname)
#print Atoms
#print Ch_NBO
#print Ch_APT
#print Ch_ESP
#print Ch_Mulliken

###print '\nCharge:'
#print 'N; Atom; Mulliken; APT; ESP; NBO;'
###print 'N; Atom; Mulliken; APT; NBO; ESP;'
###for i in range(len(Atoms)) :
#    print '%03ld' % i
###    print '%s; %s; %6.3f; %6.3f; %6.3f; %6.3f;' % (M[i], Atoms[i], Ch_Mulliken[i], Ch_APT[i], Ch_NBO[i], Ch_ESP[i])

mol_prop = {
    'Atoms': Atoms,
    'TD' : { 'E0': E0, 'H': H, 'S': S, 'G': G },
    'Wiberg': mw,
    'Freq': mfreq,
    'Ch_NBO': Ch_NBO,
    'Ch_Mulliken': Ch_NBO,
    'Ch_APT': Ch_NBO,
    'Ch_ESP': Ch_NBO,
}

print json.dumps(mol_prop, indent=4)

with open(outname, "w") as out_file:
    json.dump(mol_prop, out_file, indent=4)
