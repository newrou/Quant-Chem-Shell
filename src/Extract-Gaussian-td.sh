#!/bin/bash

list=`ls *.log`
for i in $list
do
  Name=`basename $i .log`
  TDName=$Name.td
  xyz=`basename $i .log`.xyz
  prop=`basename $i .log`.prop

  if [[ ! -e "$TDName" ]]
    then

    echo "Extract TD from Gaussian file: $Name"

    echo -n "E(B3LYP)= "  > $TDName
    grep "SCF Done:  E([RU]B3LYP) =" <$i | tail -n 1 | awk '{print $5}' >> $TDName

    echo -n "Zero-point correction= " >> $TDName
    grep "Zero-point correction=" <$i | awk '{print $3}' >> $TDName

    echo -n "Thermal correction to Energy= " >> $TDName
    grep "Thermal correction to Energy=" <$i | awk '{print $5}' >> $TDName

    echo -n "Thermal correction to Enthalpy= " >> $TDName
    grep "Thermal correction to Enthalpy=" <$i | awk '{print $5}' >> $TDName

    echo -n "Thermal correction to Gibbs Free Energy= " >> $TDName
    grep "Thermal correction to Gibbs Free Energy=" <$i | awk '{print $7}' >> $TDName

    echo -n "Sum of electronic and zero-point Energies= " >> $TDName
    grep "Sum of electronic and zero-point Energies=" <$i | awk '{print $7}' >> $TDName

    echo -n "Sum of electronic and thermal Energies= " >> $TDName
    grep "Sum of electronic and thermal Energies=" <$i | awk '{print $7}' >> $TDName

    echo -n "Sum of electronic and thermal Enthalpies= " >> $TDName
    grep "Sum of electronic and thermal Enthalpies=" <$i | awk '{print $7}' >> $TDName

    echo -n "Sum of electronic and thermal Free Energies= " <$i >> $TDName
    grep "Sum of electronic and thermal Free Energies=" <$i | awk '{print $8}' >> $TDName

# sed -e "s/\./\,/" <$Name >$Name.csv

    fi

  if [[ ! -e "$xyz" ]]
    then
    echo "Extract xyz from Gaussian file: $Name"
    ~/g16/newzmat -ichk `basename $i .log`.chk -oxyz $xyz
    n=`wc -l < $xyz`
    sed -i -e "1 s/^/$n\n$Name\n/;" $xyz
#    ~/g16/newzmat -ichk `basename $i .log`.chk -oxyz $xyz-tmp
#    n=`wc -l < $xyz-tmp`
#    echo $n > $xyz
#    echo $Name >> $xyz
#    cat $xyz-tmp >> $xyz
#    rm $xyz-tmp
    fi

  if [[ ! -e "$prop" ]]
    then
    echo "Extract property from Gaussian file: $Name"
    ./extract-g16-prop.py $i > $prop
    fi
done
