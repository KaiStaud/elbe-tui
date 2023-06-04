#!/bin/bash
template_dir=$1
output_dir=$2
version=$3
arch=$4
defconfig=$5
src_package=$6
package_type=$8

#if [ $package_type == "kernel" ]; then
echo $defconfig
echo $version
echo $arch 

cp -r $template_dir/debian $output_dir
cd $output_dir/debian
ls -la
rename "s/6.1.27/$version/" *
rename "s/stm32/$src_package/" *
sed -i  "s/stm32mp157a-dk1_defconfig/$defconfig/g" rules
sed -i -e "s/stm32/$src_package/g" control
sed -i -e "s/stm32/$src_package/g" changelog
sed -i -e "s/6.1.27/$version/g" control
sed -i -e "s/5.10.178/$version/g" rules

#elif [ $package_type == "bootloader" ]; then
#echo "not implemented"
#elif [ $package_type == "application" ]; then
#echo "Please run debmake instead"
#fi


