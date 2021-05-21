# extensions/tanzu-sql

A brief description of extensions/tanzu-sql goes here.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly. Sizing includes additional [dependencies](#dependencies).

The following table show the operator sizing requirements.

| vCPU | Memory | Storage |
| --- | --- | --- |
| 500m | 1Gi | N/A - no persistent storage required |


## Variables

### Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name                           	| Description                                       	                      | Default Value                        	| Variable Type 	| Required 	|
|-----------------------------------------	|---------------------------------------------------------------------------|--------------------------------------	|---------------	|----------	|
| tanzu_sql.namespace                     	| Namespace for extensions/tanzu-sql components                           	| "tanzu-sql"                          	| string        	| yes      	|
| tanzu_sql.staging_dir                   	| Local directory to write the staging manfiests to 	                      | "/tmp/staging/extensions/tanzu-sql" 	| string        	| yes      	|
| tanzu_sql.demo.storage_class            	| Storage class to use for demo postgresql instance                       	| "gp2-test"                           	| string        	| yes      	|
| tanzu_sql.registry.project.project_name 	| Harbor Registry name to use                       	                      | sql                                  	| string        	| yes      	|
| tanzu_sql.registry.source_url           	| Source URL to pull vmware supported instances.    	                      | dev.registry.pivotal.io              	| string        	| yes      	|
| tanzu_sql.workload_tenancy.enabled        | Whether to use the `workload-tenancy` module to provide custom namespaces | false                                 | boolean         | yes       |

```
At the time of this writing, tanzu sql for k8s is still under Beta development.
The access to the registry is restricted and request for access can be
granted by tanzu postgres PM ( lijing@vmware.com ).
```


### Additional Variables

The following variables must be set for proper operation of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name         	| Description                                                                  	| Default Value 	| Variable Type 	| Required 	|
|-----------------------	|------------------------------------------------------------------------------	|---------------	|---------------	|----------	|
| tanzu_kubectl_context 	| Name of context to use for connection to Kubernetes                          	| -             	| string        	| yes      	|
| tanzu_sql_username    	| credentials for source registry (refer to : tanzu_sql.registry.source_url )  	| -             	| string        	| yes      	|
| tanzu_sql_password    	| credentials for source registry (refer to : tanzu_sql.registry.source_url )  	| -             	| string        	| yes      	|

## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* storage
* security
* ingress
* container-registry


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=extensions/tanzu-sql make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=extensions/tanzu-sql make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

Brief description of what the demonstration does here.

In order to demonstrate the role:

```bash
ROLE=extensions/tanzu-sql make demo.role
```

Sample output:

```bash

kubectl get pods -n tanzu-sql | grep rpk-postgres

NAME                                    READY   STATUS      RESTARTS   AGE
rpk-postgres-small-0                    1/1     Running     0          24m
rpk-postgres-small-monitor-0            1/1     Running     0          24m

```

Accessing the service:

```bash
## run the following command to retrieve the database credentials

kubectl  get secret rpk-postgres-small-db-secret -n tanzu-sql -o yaml

apiVersion: v1
data:
  dbname: aXZ5LXBvc3RncmVzLWRlbW8=  ## rpk-postgres-demo
  instancename: aXZ5LXBvc3RncmVzLXNtYWxs ## rpk-postgres-small
  namespace: dGFuenUtc3Fs ## tanzu-sql
  password: YWRocXV3OEhJbDJYSEpJcUlvaWNZRGZQSTYzVGpI ## auto generated
  username: aXZ5 ## rpk
kind: Secret


```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=extensions/tanzu-sql make clean.role
```


## Author(s)
[Robin Foe](mailto:rfoe@vmware.com)
