# Kubernetes Operator with Kubebuilder - Complete Notes

---

# Goal

We are building a **custom Kubernetes Operator** that acts on our behalf and communicates with **AWS** to create, update, and delete **EC2 instances**.

Instead of manually creating EC2 instances from the AWS Console, users simply create a Kubernetes Custom Resource (CR), and the Operator handles the AWS API calls.

Flow:

```text
kubectl apply ec2instance.yaml
            │
            ▼
Kubernetes API Server
            │
            ▼
EC2Instance Controller (Operator)
            │
            ▼
AWS SDK
            │
            ▼
AWS EC2 Instance Created
```

---

# What is an Operator?

An **Operator** is a Kubernetes application that continuously watches Kubernetes resources and ensures the **actual state matches the desired state**.

For example:

Desired State:

```yaml
I want one EC2 instance.
```

Actual State:

```text
No EC2 instance exists.
```

The Operator notices the difference and creates the EC2 instance automatically.

This continuous checking is called the **Reconciliation Loop**.

---

# Kubebuilder

Kubebuilder is a framework that generates the boilerplate (scaffolding) required to build Kubernetes Operators.

It creates:

- API definitions
- Controller
- CRD
- Webhooks
- Project structure

You only write the business logic.

---

# Commands

## 1. Create a Project Folder

```bash
mkdir operator
```

### What it does

Creates a new directory named **operator**.

---

## 2. Move into the Directory

```bash
cd operator
```

### What it does

Moves into the project folder where the operator will be created.

---

## 3. Initialize Kubebuilder Project

```bash
kubebuilder init --domain cloud.com --repo github.com/shkatara/operator
```

### What it does

Initializes a Kubebuilder project.

It creates the complete project structure and configuration files.

### Meaning of each option

`--domain`

Defines the API group domain.

Example:

```
cloud.com
```

This later becomes

```
compute.cloud.com
```

inside the CRD.

---

`--repo`

Specifies the Go module path.

Example

```
github.com/shkatara/operator
```

This is used by Go modules for imports.

---

## 4. Create the API

```bash
kubebuilder create api --group compute --version v1 --kind EC2Instance
```

### What it does

**This command DOES NOT create an EC2 instance.**

Instead, it creates the framework (boilerplate/scaffolding) needed for the Operator to understand a new Kubernetes resource named **EC2Instance**.

It generates:

- Go Types
- Controller
- CRD
- RBAC
- Webhook configuration (optional)

After this command, **you still have to write the business logic**.

---

Flow:

```text
kubebuilder create api
           │
           ▼
 Generates Boilerplate
           │
 ┌─────────┼─────────┐
 ▼         ▼         ▼
Go Types  Controller  CRD
           │
           ▼
      You write
   business logic
           │
           ▼
kubectl apply ec2.yaml
           │
           ▼
 Kubernetes API Server
           │
           ▼
EC2Instance Controller
           │
           ▼
 AWS SDK
           │
           ▼
 AWS EC2 Instance Created
```

---

# Project Structure

Important directories:

```text
api/
internal/
cmd/
config/
```

---

# api/v1/

Contains the **Custom Resource Definition (CRD) types**.

Important file:

```text
api/v1/ec2instance_types.go
```

This defines what users are allowed to specify inside the YAML.

Example:

```go
type EC2InstanceSpec struct {
    AmiId  string
    SSHKey string
    Type   string
}
```

Later, users can write:

```yaml
spec:
  amiId: mydummyamiId
  sshKey: my-keypair
  type: t3.micro
```

The Go struct defines the schema for the `spec` section.

---

## Optional Fields

If a field is optional, add:

```go
omitempty
```

Example

```go
SSHKey string `json:"sshKey,omitempty"`
```

This means the field may or may not be present.

Without `omitempty`, the field is expected to always exist.

---

# Creating Nested Objects

Go doesn't recognize custom nested structures automatically.

Example:

```yaml
storage:
  standard:
    size: 10Gi
  fast:
    size: 50Gi
```

To support this, create another Go struct.

Example:

```go
type Storage struct {
    Standard Volume
    Fast Volume
}

type Volume struct {
    Size string
}
```

Then include it inside:

```go
type EC2InstanceSpec struct {
    Storage Storage
}
```

---

# Minimum YAML Required

Every Kubernetes resource needs some mandatory fields.

Example:

```yaml
apiVersion: compute.cloud.com/v1
kind: EC2Instance

metadata:
  name: myinstance
  namespace: default

spec:
  amiId: mydummyamiId
  sshKey: my-keypair
  type: t3.micro

  storage:
    standard:
      size: 10Gi

    fast:
      size: 50Gi

status:
```

---

# Required Kubernetes Metadata

## TypeMeta

Contains

```yaml
apiVersion
kind
```

Example

```yaml
apiVersion: compute.cloud.com/v1
kind: EC2Instance
```

---

## ObjectMeta

Contains

```yaml
metadata:
```

such as

- name
- namespace
- labels
- annotations

Example

```yaml
metadata:
  name: myinstance
  namespace: default
```

---

# Controller

Location:

```text
internal/controller/ec2instance_controller.go
```

This is the **heart of the Operator**.

This file contains the **Reconcile()** function.

Its responsibility is to ensure

```text
Current State == Desired State
```

Example

Desired State:

```yaml
spec:
  type: t3.micro
```

Actual State

```
No EC2 instance exists
```

The controller detects the mismatch and creates the EC2 instance through the AWS SDK.

This process repeats continuously.

This continuous execution is called the **Reconciliation Loop**.

---

# Reconciliation Loop

Every Kubernetes controller continuously performs the following steps:

```text
Read Desired State
        │
        ▼
Read Actual State
        │
        ▼
Compare Both
        │
        ▼
Different?
        │
      Yes
        │
        ▼
Take Action
        │
        ▼
Wait for Next Event
```

The controller keeps running as long as it is alive.

---

# cmd/main.go

When the Operator starts, the **first file that runs is**

```text
cmd/main.go
```

Its responsibilities include:

- Creating the Controller Manager
- Registering controllers
- Registering webhooks
- Starting metrics server
- Starting certificate watcher
- Starting leader election (if enabled)

---

# Controller Manager

The Manager is responsible for managing everything inside the Operator.

It can contain:

- One or more controllers
- Webhook server
- Metrics server
- Certificate watcher
- Cache
- Client

Diagram:

```text
Manager
│
├── EC2 Controller
├── Other Controllers
├── Webhook Server
├── Metrics Server
├── Cert Watcher
└── Cache
```

---

# Leader Election

Leader Election ensures that when multiple copies (replicas) of the Operator are running, **only one instance actively performs reconciliation**, while the others stay on standby.

This prevents multiple controllers from trying to reconcile the same resource simultaneously.

Leader election can be enabled in `cmd/main.go`.

---

# Mutating Webhook

A Mutating Webhook intercepts Kubernetes API requests **before the object is stored** and can modify the resource.

Example:

User submits:

```yaml
spec:
  type: t3.micro
```

Webhook automatically adds:

```yaml
storage:
  standard:
    size: 10Gi
```

before Kubernetes saves the object.

---

# Why HTTPS?

Communication between the Kubernetes API Server and the webhook happens over **HTTPS**, not HTTP.

Therefore, the webhook must have a valid TLS certificate.

---

# Cert Manager

Cert Manager automatically:

- Issues certificates
- Renews certificates
- Rotates certificates (typically every 90 days)

Without Cert Manager, certificate management would be manual.

---

# Problem with Certificate Rotation

Suppose Cert Manager rotates the certificate.

The webhook server now has a new certificate, but the running controller is still using the old one.

Without any additional mechanism, you would need to restart the controller so it loads the new certificate.

That causes downtime.

---

# Cert Watcher

Cert Watcher solves this problem.

It watches the certificate files.

Whenever the certificate changes, it automatically reloads the new certificate **without restarting the controller**.

Benefits:

- No downtime
- No manual restart
- Seamless certificate rotation

---

# Metrics (Prometheus)

Kubebuilder exposes a **Prometheus metrics endpoint**.

It provides metrics such as:

- Number of reconciliations executed
- Number of successful reconciliations
- Number of failed reconciliations
- Controller runtime statistics

Prometheus scrapes this endpoint and stores the metrics for monitoring.

---

# Scaffolding

Scaffolding means **automatically generated boilerplate code**.

Kubebuilder scaffolds:

- API types
- Controller
- RBAC
- Webhooks
- Configuration files
- Project structure

You only implement the business logic.

---

# Viewing Available APIs

Command:

```bash
kubectl api-resources
```

### What it does

Lists all Kubernetes resources available in the cluster.

Example output:

```
pods
deployments
services
configmaps
EC2Instances
```

This is useful to verify that your Custom Resource has been installed successfully.

---

# Updating the CRD After Spec Changes

Whenever you modify

```text
api/v1/ec2instance_types.go
```

the generated CRD becomes outdated.

You must regenerate it.

Command:

```bash
make manifests
```

### What it does

Regenerates the CRD YAML files based on your updated Go structs.

It updates:

```text
config/crd/bases/compute.cloud.com_ec2instance.yaml
```

Run this **every time you modify the Spec or Status structs**.

---

# Installing the Updated CRD

After regenerating the manifests, update the Kubernetes cluster.

Command:

```bash
make install
```

### What it does

Installs or updates the CRD inside the Kubernetes cluster.

Flow:

```text
Modify Spec
      │
      ▼
make manifests
      │
      ▼
Update CRD YAML
      │
      ▼
make install
      │
      ▼
Cluster Updated
```

---

# What is an EC2 Instance?

An **EC2 (Elastic Compute Cloud) Instance** is a virtual machine running inside AWS.

Instead of buying a physical server, you rent a virtual server from AWS.

It behaves like a normal computer with:

- CPU
- RAM
- Storage
- Operating System
- Network Interface
- Public and Private IP addresses

Example:

```text
Ubuntu Server
CPU: 2 vCPUs
RAM: 4 GB
Disk: 50 GB SSD
Public IP: 54.xx.xx.xx
```

---

# Why Use EC2?

EC2 instances are commonly used to:

- Host websites
- Run backend services
- Run databases
- Train machine learning models
- Run Docker containers
- Host Kubernetes worker nodes
- Execute scripts and batch jobs

Example:

```text
You build a Node.js application.

Instead of running it on your laptop:

Launch an EC2 instance
SSH into it
Install Node.js
Run the application

Users can now access your application over the internet.
```

---

# What Happens When an EC2 Instance is Created?

When you launch an EC2 instance, AWS finds an available physical server in one of its data centers and creates a virtual machine for you.

```text
AWS Physical Server
+--------------------------------------+
|                                      |
|  VM 1 (Ubuntu EC2)                   |
|                                      |
|  VM 2 (Amazon Linux EC2)             |
|                                      |
|  VM 3 (Windows EC2)                  |
|                                      |
+--------------------------------------+
```

Multiple EC2 instances can run on the same physical machine.

---

# Important EC2 Properties

## 1. AMI (Amazon Machine Image)

The template used to create the virtual machine.

Examples:

- Ubuntu 24.04
- Amazon Linux
- Red Hat
- Windows Server

Example in the CRD:

```yaml
amiId: ami-0123456789abcdef
```

---

## 2. Instance Type

Defines the hardware configuration of the virtual machine.

| Instance Type | CPU | RAM |
|---------------|-----|-----|
| t2.micro | 1 vCPU | 1 GB |
| t3.micro | 2 vCPU | 1 GB |
| t3.small | 2 vCPU | 2 GB |
| t3.medium | 2 vCPU | 4 GB |
| m5.large | 2 vCPU | 8 GB |

Example:

```yaml
type: t3.micro
```

---

## 3. SSH Key Pair

Used to securely log into the EC2 instance.

Example:

```yaml
sshKey: my-keypair
```

Connect using:

```bash
ssh -i my-keypair.pem ubuntu@54.xx.xx.xx
```

---

## 4. Storage (EBS Volumes)

Persistent disks attached to the EC2 instance.

Examples:

- 20 GB SSD
- 100 GB SSD
- 500 GB HDD

---

## 5. Security Group

Acts as a virtual firewall.

Controls which network ports are open.

Example:

```text
22   -> SSH
80   -> HTTP
443  -> HTTPS
```

---

## 6. Public IP

If enabled, AWS assigns a public IP address.

Example:

```text
18.192.44.121
```

Users can then access the server over the internet.

Example:

```text
http://18.192.44.121
```

---

# Complete Example

Suppose you create an EC2 instance with:

```text
AMI        : Ubuntu 24.04
Type       : t3.micro
Storage    : 20 GB SSD
SSH Key    : developer-key
```

AWS creates a virtual machine:

```text
Operating System : Ubuntu 24.04
CPU              : 2 vCPUs
RAM              : 1 GB
Disk             : 20 GB SSD
Public IP        : 54.210.11.90
```

You connect:

```bash
ssh -i developer-key.pem ubuntu@54.210.11.90
```

Install software:

```bash
sudo apt update
sudo apt install nginx
```

Start NGINX:

```bash
sudo systemctl start nginx
```

Visit:

```text
http://54.210.11.90
```

You should see the default NGINX welcome page.

---

# Overall Workflow

```text
mkdir operator
        │
        ▼
cd operator
        │
        ▼
kubebuilder init
        │
        ▼
kubebuilder create api
        │
        ▼
Generated:
    ├── API
    ├── Controller
    ├── CRD
    └── Boilerplate
        │
        ▼
Define Spec in ec2instance_types.go
        │
        ▼
Implement Reconcile() in ec2instance_controller.go
        │
        ▼
make manifests
        │
        ▼
Update CRD YAML
        │
        ▼
make install
        │
        ▼
kubectl apply ec2instance.yaml
        │
        ▼
Kubernetes API Server
        │
        ▼
EC2Instance Controller
        │
        ▼
AWS SDK
        │
        ▼
AWS Creates EC2 Instance
```