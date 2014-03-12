#!/bin/bash
# Bash Menu Script Example

PS3='Please enter your choice: '
echo "[OPTIONS]"
options=("Image list" "Container List" "Run a container" "Remove an image" "Kill a container" "Quit")
select opt in "${options[@]}"
do
    case $opt in
    
		"Image list")
			clear
			echo -e "Here is the list of the available images : \n*****"   
			sudo docker images
			echo -e "*****\n\n[OPTIONS]\n1) Image list\n2) Container list \n3) Run a container \n4) Remove an image \n5) Kill a container \n6) Quit"
			;;
    
        "Container List")
			clear
			echo -e "Here is the list of the running containers : \n*****"
            sudo docker ps
            echo -e "*****\n\n[OPTIONS]\n1) Image list\n2) Container list \n3) Run a container \n4) Remove an image \n5) Kill a container \n6) Quit"           
            ;;
            
         "Run a container")
			clear
			echo -e "Here is the list of the available images : \n*****"   
			sudo docker images
			echo -e "*****\n\nWhich container do you want to run ?"
			read nom
			sudo docker run -t -d $nom
			echo -e "Container is running \n"
			echo -e "*****\n\n[OPTIONS]\n1) Image list\n2) Container list \n3) Run a container \n4) Remove an image \n5) Kill a container \n6) Quit"
            ;;
            
         "Remove an image")
			clear
			echo -e "Here is the list of the available images : \n*****"
			sudo docker images
			echo -e "*****\n\nWhich image do you want to remove ?"
			read nom
			sudo docker rmi $nom
			echo -e "Image removed. \n"
			echo -e "*****\n\n[OPTIONS]\n1) Image list\n2) Container list \n3) Run a container \n4) Remove an image \n5) Kill a container \n6) Quit"  
			;;
			
		 "Kill a container")
			clear
			echo -e "Here is the list of the running containers : \n*****"   
			sudo docker ps
			echo -e "*****\n\nWhich container do you want to kill ?"
			read nom
			sudo docker kill $nom
			echo -e "Container is killed \n"
			echo -e "*****\n\n[OPTIONS]\n1) Image list\n2) Container list \n3) Run a container \n4) Remove an image \n5) Kill a container \n6) Quit"		
			;;
			
		 "Inspect")
			;;
            
        "Quit")
            break
            ;;
        *) echo invalid option;;
    
    esac
done
