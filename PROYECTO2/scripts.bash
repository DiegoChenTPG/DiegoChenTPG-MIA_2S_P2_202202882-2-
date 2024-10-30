mkdisk -size=5 -unit=M -fit=WF -path=discos/DiscoLab.mia
mkdisk -size=5 -unit=M -fit=WF -path="discos/DiscoLab1.mia"

fdisk -size=1 -type=P -unit=M -fit=BF -name="Particion1" -path="discos/DiscoLab.mia"
fdisk -size=2 -type=P -unit=M -fit=WF -name="Particion2" -path="discos/DiscoLab.mia"
fdisk -size=1 -type=P -unit=M -fit=BF -name="Particion1" -path="discos/DiscoLab1.mia"

mount -name="Particion1" -path="discos/DiscoLab.mia"
mount -name="Particion2" -path="discos/DiscoLab.mia"
mount -name="Particion1" -path="discos/DiscoLab1.mia"

mkfs -id=820A
mkfs -id=821A -type=full

rmdisk -path="discos/DiscoLab.mia"
rmdisk -path="discos/DiscoLab1.mia"



rep -id=820A -path="salidas/reporte_mbr.png" -name=mbr
rep -id=820A -path="salidas/reporte_inde.png" -name=inode
rep -id=820A -path="salidas/reporte_bm_inode.txt" -name=bm_inode





mkdisk -size=5 -unit=M -fit=WF -path=discos/DiscoLab.mia


fdisk -size=1 -type=P -unit=M -fit=BF -name="Particion1" -path="discos/DiscoLab.mia"
fdisk -size=2 -type=P -unit=M -fit=WF -name="Particion2" -path="discos/DiscoLab.mia"

mount -name="Particion1" -path="discos/DiscoLab.mia"
mount -name="Particion2" -path="discos/DiscoLab.mia"


mkfs -id=820A
mkfs -id=821A -type=full

rmdisk -path="discos/DiscoLab.mia"



rep -id=820A -path="salidas/reporte_mbr.png" -name=mbr
rep -id=820A -path="salidas/reporte_inde.png" -name=inode
rep -id=820A -path="salidas/reporte_bm_inode.txt" -name=bm_inode




mkdisk -size=3000 -unit=K -fit=WF -path=discos/DiscoLab.mia


fdisk -size=300 -type=P -unit=K -fit=WF -name="Particion1" -path="discos/DiscoLab.mia"

fdisk -size=300 -type=E -unit=K -fit=WF -name="Particion2" -path="discos/DiscoLab.mia"

fdisk -size=100 -type=L -unit=K -fit=BF -name="Particion3" -path="discos/DiscoLab.mia"


mount -name="Particion1" -path="discos/DiscoLab.mia"
mount -name="Particion2" -path="discos/DiscoLab.mia"
mount -name="Particion3" -path="discos/DiscoLab.mia"









#Calificacion Proyecto 1
#2S 2024
 
#----------------- 1. MKDISK  -----------------


#----------------- MKDISK CON ERROR -----------------
# ERROR PARAMETROS
mkdisk -param=x -size=30 -path=/home/drct-tpg/Calificacion_MIA/Discos/DiscoN.mia


#----------------- CREACION DE DISCOS -----------------
# ERROR PARAMETROS
mkdisk -tamaño=3000 -unit=K -path=/home/drct-tpg/Calificacion_MIA/Discos/DiscoN.mia
# 50M A
Mkdisk -size=50 -unit=M -fit=FF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia
# 50M B
Mkdisk -unit=k -size=51200 -fit=BF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco2.mia
# 13M C
mkDisk -size=13 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia
# 50M D
mkdisk -size=51200 -unit=K -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco4.mia
# 20M E
mkDisk -size=20 -unit=M -fit=WF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia
# 50M F X
Mkdisk -size=50 -unit=M -fit=FF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco6.mia
# 50M G X
Mkdisk -size=50 -unit=M -fit=FF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco7.mia
# 50M H X
mkdisk -size=51200 -unit=K -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco8.mia
# 50M I X
mkdisk -size=51200 -unit=K -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco9.mia
# 50M J X
mkdisk -size=51200 -unit=K -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco10.mia


#-----------------2. RMDISK-----------------
#ERROR DISCO NO EXISTE
rmdisk -path=/home/drct-tpg/Calificacion_MIA/Discos/DiscoN.mia
# BORRANDO DISCO
rmdisk -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco6.mia
# BORRANDO DISCO
rmdisk -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco7.mia
# BORRANDO DISCO
rmdisk -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco8.mia
# BORRANDO DISCO
rmdisk -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco9.mia
# BORRANDO DISCO
rmdisk -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco10.mia


#-----------------3. FDISK-----------------
#-----------------CREACION DE PARTICIONES-----------------
#DISCO 1
# ERROR RUTA NO ENCONTRADA
fdisk -type=P -unit=b -name=PartErr -size=10485760 -path=/home/drct-tpg/Calificacion_MIA/Discos/DiscoN.mia -fit=BF 
# PRIMARIA 10M
fdisk -type=P -unit=b -name=Part11 -size=10485760 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=k -name=Part12 -size=10240 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=M -name=Part13 -size=10 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=b -name=Part14 -size=10485760 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
#ERR LMITE PARTICION PRIMARIA
fdisk -type=P -unit=b -name=PartErr -size=10485760 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF


# LIBRE DISCO 1: 50-4*10 = 10 -> 20%


#DISCO 3
# ERROR FALTA ESPACIO
fdisk -type=P -unit=m -name=PartErr -size=20 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia
#4M
fdisk -type=P -unit=m -name=Part31 -size=4 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia
#4M
fdisk -type=P -unit=m -name=Part32 -size=4 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia
#1M
fdisk -type=P -unit=m -name=Part33 -size=1 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia


#LIBRE DISCO 3: 13-9= 4 -> 30.77%


#DISCO 5
# 5MB
fdisk -type=E -unit=k -name=Part51 -size=5120 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part52 -size=1024 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 5MB
fdisk -type=P -unit=k -name=Part53 -size=5120 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part54 -size=1024 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part55 -size=1024 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=k -name=Part56 -size=1024 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -fit=BF


# LIBRE DISCO 5: 20-10 = 5 -> 50%
# LIBRE EXTENDIDA 2: 5-4 = 1 -> 20% (por los EBR deberia ser menos)
#FUNCIONA HASTA ACA

#-----------------MOUNT-----------------
#-----------------MONTAR PARTICIONES-----------------
#DISCO 1
#821A 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -name=Part11
#822A 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -name=Part12
#ERROR PARTICION YA MONTADA
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -name=Part11


#DISCO 3
#ERROR PARTCION NO EXISTE
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia -name=Part0
#821B 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia -name=Part31
#822B 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco3.mia -name=Part32


#DISCO 5
#821C 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco5.mia -name=Part53


#-----------------REPORTES PARTE 1-----------------
#DISCO 1
#ERROR ID NO ENCONTRADO
rep -id=A821 -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_rE.jpg -name=mbr
#REPORTE DISK
rep -id=821A -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_r1_disk.jpg -name=disk
#REPORTE MBR 
rep -id=821A -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_r2_mbr.jpg -name=mbr


#DISCO 3
#ERROR ID NO ENCONTRADO
rep -id=823B -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_rE_mbr.jpg -name=mbr
#REPORTE DISK
rep -id=821B -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_r3_disk.jpg -name=disk
#REPORTE MBR
rep -id=822B -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_r4_disk.jpg -name=mbr


#DISCO 5
#ERROR ID NO ENCONTRADO
rep -id=IDx -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_rE_mbr.jpg -name=mbr
#REPORTE DISK
rep -id=821C -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_r5_disk.jpg -name=disk
#REPORTE MBR
rep -id=821C -Path=/home/drct-tpg/Calificacion_MIA/Reportes/p1_r6_mbr.jpg -name=mbr


#-----------------5. MKFS-----------------
mkfs -type=full -id=821A


#-----------------PARTE 3-----------------
#-----------------7. LOGIN-----------------
login -user=root -pass=123 -id=821A
#ERROR SESION INICIADA
login -user=root -pass=123 -id=821A


#-----------------9. MKGRP-----------------
mkgrp -name=usuarios
mkgrp -name=adm
mkgrp -name=mail
mkgrp -name=news
mkgrp -name=sys
#ERROR YA EXISTE EL GRUPO
mkgrp -name=sys













# esto es lo que vamos probando ahora

#Calificacion Proyecto 1
#2S 2024


#----------------- 1. MKDISK  -----------------
# 50M A
Mkdisk -size=50 -unit=M -fit=FF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia

# PRIMARIA 10M
fdisk -type=P -unit=b -name=Part11 -size=10485760 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=k -name=Part12 -size=10240 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=M -name=Part13 -size=10 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF


#DISCO 1
#821A 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -name=Part11



#-----------------5. MKFS-----------------
mkfs -type=full -id=821A


#-----------------9. MKGRP-----------------
mkgrp -name=usuarios
mkgrp -name=adm
mkgrp -name=mail
mkgrp -name=news
mkgrp -name=sys



# =================================================================



#-----------------11. MKUSR-----------------------------------------------
mkusr -user="usuario1" -pass=password -grp=root
mkusr -user="user1" -pass=abc -grp=usuarios
mkusr -user="user2" -pass=abc -grp=usuarios
#ERROR EL USUARIO YA EXISTE
mkusr -user="user2" -pass=abc -grp=usuarios
#ERROR GRUPO NO EXISTE
mkusr -user="user3" -pass=abc -grp=system







# esto es lo que vamos probando ahora

#Calificacion Proyecto 1
#2S 2024


#----------------- 1. MKDISK  -----------------
# 50M A
Mkdisk -size=50 -unit=M -fit=FF -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia

# PRIMARIA 10M
fdisk -type=P -unit=b -name=Part11 -size=10485760 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=k -name=Part12 -size=10240 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=M -name=Part13 -size=10 -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -fit=BF


#DISCO 1
#821A 
mount -path=/home/drct-tpg/Calificacion_MIA/Discos/Disco1.mia -name=Part11



#-----------------5. MKFS-----------------
mkfs -type=full -id=821A





#-----------------7. LOGIN-----------------
login -user=root -pass=123 -id=821A
#ERROR SESION INICIADA
login -user=root -pass=123 -id=821A


#-----------------9. MKGRP-----------------
mkgrp -name=users
#ERROR
mkgrp -name=users
# =================================================================

#-----------------10. RMGR-----------------
#ERROR
rmgrp -name=mail
#rmgrp -name=users


#-----------------11. MKUSR-----------------------------------------------
mkusr -user="usr1" -pass=pas -grp=root
#ERROR
mkusr -user="usr1" -pass=pass -grp=root


#-----------------12. RMUSR-----------------
#ERROR
rmusr -user=user2

#se elimina
#rmusr -user=usr1

#-----------------13. CHGRP-----------------
chgrp -user=usr1 -grp=users




