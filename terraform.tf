terraform {  
  backend "gcs" {
    project = "auzom-legacy-jeremejevs"
    bucket  = "auzom-legacy-jeremejevs"
    path    = "terraform.tfstate"
  }
}

provider "google" {
  version = "~> 1.7"
  project = "auzom-legacy-jeremejevs"
  region  = "us-east1"
}

# Be advised, Terraform's state is stored in this bucket
resource "google_storage_bucket" "auzom-legacy-jeremejevs" {
  name          = "auzom-legacy-jeremejevs"
  storage_class = "REGIONAL"
  location      = "us-east1"

  versioning {
    enabled = true
  }
}

resource "google_compute_disk" "data" {
  name = "data"
  type = "pd-ssd"
  zone = "us-east1-d"
  size = 10
}

resource "google_compute_network" "main" {
  name = "main"
}

resource "google_compute_firewall" "main" {
  name    = "main"
  network = "${google_compute_network.main.self_link}"

  allow {
    protocol = "icmp" # For IPv6
  }

  allow {
    protocol = "tcp"
    ports    = ["22", "80", "443"]
  }
}

resource "google_compute_instance" "main" {
  name         = "main"
  machine_type = "f1-micro"
  zone         = "us-east1-d"

  boot_disk {
    initialize_params {
      image = "cos-cloud/cos-stable"
    }
  }

  attached_disk {
    source      = "${google_compute_disk.data.self_link}"
    device_name = "data"
  }

  network_interface {
    network = "${google_compute_network.main.self_link}"
    access_config {}
  }

  service_account {
    scopes = ["storage-ro"] # For GCR
  }

  metadata_startup_script = <<EOF
    mkdir -p /mnt/disks/data &&
    mount /dev/disk/by-id/google-data /mnt/disks/data &&
    su - user -c 'cd ~ && . scripts/docker-compose up -d'
  EOF
}

variable "cloudflare_email" {}
variable "cloudflare_token" {}

provider "cloudflare" {
  version = "~> 0.1"
  email   = "${var.cloudflare_email}"
  token   = "${var.cloudflare_token}"
}

resource "cloudflare_record" "frontend" {
  domain = "auzom.gg"
  name   = "legacy"
  type   = "A"
  value  = "${google_compute_instance.main.network_interface.0.access_config.0.assigned_nat_ip}"
  proxied = true
}

resource "cloudflare_record" "backend" {
  domain = "auzom.gg"
  name   = "api.legacy"
  type   = "A"
  value  = "${google_compute_instance.main.network_interface.0.access_config.0.assigned_nat_ip}"
}
