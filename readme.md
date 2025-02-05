# Sqlite Backup

A lightweight tool that automates the backup process for an SQLite database. It performs the following tasks:

* Creates a backup of the specified SQLite database.
* Uploads the backup to Azure Blob Storage.
* Manages backup retention by keeping only a predefined number of backups.

## Configuration

Set the following environment variables to authenticate and store backups in Azure Blob Storage:

- `AZ_ACCOUNT` – Your Azure Storage account name.
- `AZ_ACCOUNT_KEY` – Your Azure Storage account key.
- `AZ_CONTAINER` – The Azure Blob Storage container to store backups.

## Usage

Run the tool with the following arguments:

```sh
sqlitebc <database_file_path> <azure_directory> <storage_capacity>
```

- `<database_file_path>` – Path to the SQLite database file.
- `<azure_directory>` – Directory in Azure Blob Storage where backups will be stored.
- `<storage_capacity>` - Number of files to keep in blob storage. Optional, 5 by default. Must be `int`.

## Example

```sh
export AZ_ACCOUNT="myaccount"
export AZ_ACCOUNT_KEY="myaccountkey"
export AZ_CONTAINER="mycontainer"

sqlitebc /path/to/database.db directory 5
```

This will:

* Create a backup of `/path/to/database.db`.  
* Upload it to `backups/sqlite/` in the specified Azure Blob Storage container.  
* Keep only the latest 5 backups in storage.
