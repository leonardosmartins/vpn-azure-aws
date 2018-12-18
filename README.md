# VPN-AZURE-AWS
 
 Este projeto tem a finalidade de demonstrar como criar uma VPN entre uma VM hospedada na AWS com outra VM hospedada na Azure e criar um ambiente de demonstração de funcionamento dessa comunicação.

## Descrição do ambiente de demonstração

 Para o ambiente de demonstração iremos construir um cenário onde instalaremos e configuraremos uma instância MySQL em uma das VMs. Nessa instância do MySQL criaremos um Database e um usuário com permissões completas nesse Database criado. Em seguida será criado uma tabela `filme` com um campo `FILME_ID` e outro campo `FILME_NAME`.
 
 Na outra VM, com a utilização de um script ansible, iremos instalar a ferramnte Docker e iniciaremos um ambiente de Docker Swarm. Por fim criaremos uma imagem docker a partir do `Dockerfile` presente neste repositório que irá conter a aplicação `connection.go`, também presente neste repositório, a qual inicia uma conexão com o banco MySQL e insere um valor, passado pelo script ansible, na tabela `filme` criada no nosso MySQL.  
 
# VPN

Para realizar a comunicação entre a VM hospedada na AWS com a VM hospedada na Azure, estabeleceremos uma conexão VPN utilizando uma ferramenta chamada Strongswan.

### Configurando Ambiente

Primeiro vamos listar uma sequência de recursos que devem ser criados e configurados inicialmente em cada uma das nuvens

#### Azure

 - Criar Resource Group
 - Criar Virtual Network
 - Criar Virtual Machine (Ubuntu)
 - Liberar portas 500 e 4500 com protocolo UDP na Virtual Machine criada

#### AWS

 - Criar VPC
 - Criar Subnet
 - Criar Internet Gateway
 - Associar Internet Gateway criado na VPC
 - Criar Route Table
 - Na Route Table criada adicionar a rota 0.0.0.0/0 com Target no Internet Gateway
 - Na Route Table criada associar a Subnet criada anteriormente
 - Criar EC2 (Ubuntu)
 - Liberar no Target Group da EC2 criada as portas 500 e 4500 com protocolo UDP
 - Na Network da EC2 criada desabilitar Source/Destination Check
 - Criar Elastic IP
 - Associar Elastic IP na EC2 criada

### Instalando e configurando Strongswan

Essas configurações devem ser feitas nas duas VMs, utilizando a informações adequadas em cada VM.

Atualizando pacotes:
```sh
$ sudo apt-get update
```

Instalando Strongswan:
```sh
$ sudo apt-get install strongswan
```

Editar o arquivo `ipsec.conf`. Usar template `ipsec-template.conf` preenchendo os dados necessários:
```sh
$ sudo nano /etc/ipsec.conf
```

Editar o arquivo `ipsec.secrets`. Usar template `ipsec-template.secrets` preenchendo os dados necessários:
```sh
$ sudo nano /etc/ipsec.secrets
```

Editar o arquivo `sysctl.conf` e descomentar a linha `net.ipv4.ip_forward = 1`:
```sh
$ sudo nano /etc/sysctl.conf
```

Habilitando alterações realizadas:
```sh
$ sudo sysctl -p /etc/sysctl.conf
```

Reiniciar:
```sh
$ sudo reboot
```

Verificando conexão estabelecida:
```sh
$ sudo ipsec status
```

# MySQL

Com a nossa VPN já estabelecida, vamos instalar e configurar uma instância MySQL e uma das VMs.

### Instalando e configurando

Instalando MySQL:
```sh
$ sudo apt-get install mysql-server
```

Configurações iniciais:
```sh
$ sudo mysql_secure_installation
```

Liberando acesso remoto:
```sh
$ sudo ufw allow mysql
```

Editar o arquivo `mysqld.cnf` e comentar a linha `bind-address = 127.0.0.1`  :
```sh
$ sudo nano /etc/mysql/mysql.conf.d/mysqld.cnf
```

Reiniciando MySQL:
```sh
$ sudo /etc/init.d/mysql restart
```

### Criando ambiente de demonstração

Conectando como root:
```sh
$ sudo mysql -u root -p
```

Criando projeto database:
```sh
mysql> CREATE DATABASE projeto;
```

Criando usuário exemplo:
```sh
mysql> CREATE USER 'exemplo'@'%' IDENTIFIED BY 'exemplopassword';
```

Configurando permissões para o usuário no database projeto:
```sh
mysql> GRANT ALL ON projeto.* to exemplo@'%' IDENTIFIED BY 'exemplopassword' WITH GRANT OPTION;
```

Recarregando permissões:
```sh
mysql> FLUSH PRIVILEGES;
```

Conectando ao database:
```sh
mysql> USE projeto;
```

Criando tabela filme:
```sh
mysql> CREATE TABLE filme (filme_id INT NOT NULL AUTO_INCREMENT, filme_name VARCHAR(50) NOT NULL, PRIMARY KEY (filme_id));
```

# Ansible

Utilizaremos a ferramenta Ansible para a automatização da configuração do nosso ambiente de Docker Swarm e o deploy da aplicação que acessará o MySQL que criamos.

Antes de rodarmos nosso script ansible, deve-se primeiro preencher com as informações necessárias no arquivo ansible/hosts

Deve-se também antes de executar o script ansible conectar na VM na qual configuraremos esse ambiente, e instalar alguns pacotes:

Instalando python:
```sh
$ sudo apt-get install python python-pip
```

Instalando docker-py:
```sh
$ sudo pip install docker-py
```

Por fim deve-se colocar o connection.go e o Dockerfile dentro dessa VM no diretório `/home/ubuntu/mysql-vpn`

Executando nosso ansible:
```sh
$ ansible-playbook -i hosts main.yml
```

Quando solicitado deve-se informar as seguintes informações:
 - db_ip = ip_privado_mysql
 - db_name = database_name
 - db_user = usuário_criado_mysql
 - db_password = password_criado_mysql
 - value = valor desejado para inserir como filme_name

# Validação

Para verificar todo nosso processo, basta entrar na VM com o MySQL, conectar ao MySQL e rodar os seguintes comando:

Conectando ao database:
```sh
mysql> USE projeto;
```

Select na tabela filme:
```sh
mysql> SELECT * FROM filme;
```
