# Application Pipeline

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md).

This [Guide](https://tanzu.vmware.com/developer/guides/ci-cd/)
demonstrates deployment of secure software pipeline with DevOps
toolchain on Kubernetes using the following tools:

* Git
* SonarQube
* Nexus Repository
* Jenkins Pipeline
  * Modular design
  * Shared Library
  * Tools Integrations
  * Deployment to kubernetes cluster with kubectl


## Default Variables

### URL Variables
| Variable Name                              	| Description                                      	| Default Value                        	| Variable Type 	| Required 	|
|--------------------------------------------	|--------------------------------------------------	|--------------------------------------	|---------------	|----------	|
| tanzu_app_pipeline_gateway                 	| domain name pointing at contour.envoy.externalIP 	| pipeline.tanzu-no.de                 	| string        	| yes      	|
| tanzu_app_pipeline_gateway_git_http        	| git HTTP URL                                     	| gitea.pipeline.tanzu-no.de           	| string        	| yes      	|
| tanzu_app_pipeline_gateway_git_ssh         	| git SSH URL                                      	| tunnel.gitea.pipeline.tanzu-no.de    	| string        	| yes      	|
| tanzu_app_pipeline_gateway_sonarqube       	| sonarqube URL                                    	| sonarqube.pipeline.tanzu-no.de       	| string        	| yes      	|
| tanzu_app_pipeline_gateway_nexus_container 	| nexus container URL                              	| container.nexus.pipeline.tanzu-no.de 	| string        	| yes      	|
| tanzu_app_pipeline_gateway_nexus_http      	| nexus proxy URL                                  	| proxy.nexus.pipeline.tanzu-no.de     	| string        	| yes      	|
| tanzu_app_pipeline_gateway_jenkins         	| jenkins URL                                      	| jenkins.pipeline.tanzu-no.de         	| string        	| yes      	|


### Tools Variables

This variable start with prefix `tanzu_app_pipeline.tools`


#### Gitea

Prefix `tanzu_app_pipeline.tools.git`

Default username : rpk

Default password : tanzu

| Variable Name                    	| Description                                                        | Default Value                                         	| Variable Type 	| Required 	|
|----------------------------------	|--------------------------------------------------------------------|-------------------------------------------------------	|---------------	|----------	|
| enable                           	| flag to enable tools deployment                                    | true                                                  	| boolean       	| yes      	|
| name                             	| name of the helm chart                                             | gitea                                                 	| string        	| yes      	|
| default.storage.class           	| Storage class                                                      | gp2-test                                               | string        	| yes       |
| default.credential.username      	| admin username                                                     | rpk                                                   	| string        	| yes      	|
| default.credential.password      	| admin password                                                     | tanzu                                          	| string        	| yes      	|
| default.credential.email         	| admin email, free style email address                              | rpk@{{ tanzu_ingress_domain }}                                    	| string        	| yes      	|
| default.repo.pipeline.clone_addr 	| repository to jenkins groovy pipeline                              | https://github.com/robinfoe/software-supply-chain.git 	| string        	| yes      	|
| default.repo.pipeline.repo_name  	| local repository name, will create if absent                       | software-supply                                       	| string        	| yes      	|
| default.repo.app.clone_addr      	| application repository location                                    | https://github.com/robinfoe/bookstore-ms.git          	| string        	| yes      	|
| default.repo.app.repo_name       	| local repository name, will create if absent                       | bookstore-ms                                          	| string        	| yes      	|


#### Sonatype Nexus

Prefix `tanzu_app_pipeline.tools.nexus`

Default username : admin

Default password : admin123

| Variable Name     	  | Description                                                                               	| Default Value  	| Variable Type 	| Required 	|
|---------------------	|-------------------------------------------------------------------------------------------	|----------------	|---------------	|----------	|
| enable            	  | flag to enable tools deployment                                                             | true           	| boolean       	| yes      	|
| name              	  | name of the helm chart                                                                    	| sonatype-nexus 	| string        	| yes      	|
| default.storage.class | default storage class                                    	                                  | gp2-test       	| string        	| yes      	|
| default.password  	  | password generated by helm installation, do not change                                    	| admin123       	| string        	| yes      	|
| default.anonymous 	  | allow anonymous access                                                                    	| true           	| boolean       	| yes      	|
| default.repo.name 	  | repository name to store application binary build with rpk pipeline                       	| rpk-snapshot   	| string        	| yes      	|


#### SonarQube

Prefix `tanzu_app_pipeline.tools.sonar`

Default username : admin

Default password : admin

| Variable Name     	      | Description                           	| Default Value  	| Variable Type 	| Required 	|
|-------------------------	|---------------------------------------	|----------------	|---------------	|----------	|
| enable            	      | flag to enable git deployment           | true           	| boolean       	| yes      	|
| name              	      | name of the helm chart                 	| sonarqube     	| string        	| yes      	|
| default.database.password | database password                       | c29uYXJQYXNz         	| string        	| yes      	|



#### Jenkins

Prefix `tanzu_app_pipeline.tools.jenkins`

Default username : admin

Default password : tanzu

| Variable Name                	| Description                        	| Default Value                                   	| Variable Type 	| Required 	|
|------------------------------	|------------------------------------	|-------------------------------------------------	|---------------	|----------	|
| enable                       	| flag to enable tools deployment    	| true                                            	| boolean       	| yes      	|
| name                         	| name of the helm chart             	| jenkins                                         	| string        	| yes      	|
| default.storage.class         | default storage class               | gp2-test       	                                  | string        	| yes      	|
| default.registry.credential 	| base64 token for dockerhub access  	| -                                               	| string        	| yes      	|
| default.credential.username 	| base64 username                   	| YWRtaW4=                                         	| string        	| yes      	|
| default.credential.password 	| base64 password                    	| aXZ5YWRtaW4=                                    	| string        	| yes      	|
| pipeline.repo                	| repository to pull pipeline script 	| {{ path to gitea }}/rpk/software-supply.git         | string        	| yes      	|
| pipeline.branch              	| repository branch                  	| rpk-1.0                                         	| string        	| yes      	|
| pipeline.files               	| groovy script to include as jobs   	| pipe-ci, pipe-full, pipe-cb-dockerfile, pipe-cd 	| array         	| yes      	|


## Dependencies

* pip
* docker
* kubectl
* git
* curl
* jq
* python-jenkins

## pip Installation

```bash
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python get-pip.py
```

## Dependencies Installation

```bash
cd {{ role_path }}
~/Library/Python/2.7/bin/pip install -r requirements.txt
~/Library/Python/2.7/bin/pip install --upgrade --user openshift
```

## Running application-pipeline Playbook

The following must be run from the root of the `reference-platform-for-kubernetes` project, and NOT from this role directory:

```bash
ansible-playbook -vv test.yaml \
        -e 'ansible_python_interpreter=./ansible-virtualenv/bin/python3' \
        -e 'rpk_role_name=application-pipeline' \
        -c local -i build/inventory.yaml
```

## Watch a Demo

[Application-Pipeline Demo](https://drive.google.com/file/d/17qFo9w_VvV45OmHbmF7x3jm8tKOYVVJc/view?usp=sharing)
