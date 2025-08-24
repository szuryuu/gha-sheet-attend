# GHA Sheet Attend - Daily Attendance via GitHub Actions

This repository contains a GitHub Actions implementation to automatically log daily attendance (e.g., for internships or work) into a Google Sheet directly from the GitHub UI.

---

### Features
- **Simple Input Form**: Fill in your daily attendance through an easy-to-use form in the GitHub "Actions" tab.
- **Full Automation**: The script automatically fetches the current date and calculates the next sequential number.
- **Automatic Table Formatting**: Every new entry is automatically formatted with table borders to keep the sheet tidy.
- **Secure**: Credentials are stored using GitHub Secrets, not hardcoded.

### Setup Instructions
Follow these steps to configure this project for your own use.

#### **Step 1: Get the Code**
- **Fork**, **Clone**, or **Use as Template** this repository to your GitHub account.

#### **Step 2: Configure Google Cloud & Service Account**
You need a "bot account" (Service Account) to allow the script to access your Google Sheet.
1.  Go to the [Google Cloud Console](https://console.cloud.google.com/) and create a **New Project**.
2.  Within the project, enable two APIs: **Google Sheets API** and **Google Drive API**.
3.  Create a **Service Account**:
    - Navigate to `APIs & Services` > `Credentials`.
    - Click `Create Credentials` > `Service account`.
    - Give it a name (e.g., `github-actions-writer`), then click `Create and Continue`.
4.  Generate a **JSON Key** for the Service Account:
    - Click on the newly created Service Account.
    - Go to the `Keys` tab > `Add Key` > `Create new key`.
    - Choose the **JSON** format and click `Create`. A `.json` file will be downloaded. **Keep this file safe; its contents are secret.**

#### **Step 3: Configure Google Sheets**
1.  Create a **new Google Sheet**.
2.  Create the headers in the first few rows, ensuring your data starts from row 5, with the columns:
    `No`, `Hari/Tanggal`, `Waktu Mulai`, `Waktu Selesai`, `Keterangan`, `Keterangan tambahan`
3.  **Share** your Sheet:
    - Open the `.json` file you downloaded, and copy the email address inside (e.g., `...gserviceaccount.com`).
    - In your Google Sheet, click `Share`, paste the email address, and grant it **Editor** access.
4.  Get the **Spreadsheet ID** and **Sheet ID (gid)**:
    - **Spreadsheet ID**: Found in the URL, e.g., `.../spreadsheets/d/`**`THIS_IS_THE_ID`**`/edit...`
    - **Sheet ID (gid)**: Found at the end of the URL, e.g., `.../edit#gid=`**`123456789`**. For the first sheet, this is usually `0`.

#### **Step 4: Configure GitHub Repository**
1.  **Add a Secret**:
    - In your GitHub repository, go to `Settings` > `Secrets and variables` > `Actions`.
    - Click `New repository secret`.
    - Name the secret: `GCP_SA_KEY`.
    - Open your `.json` key file, copy its **entire content**, and paste it into the secret's value.
2.  **Update the `main.go` Code**:
    - Open the `main.go` file.
    - Replace the constant values at the top with the IDs you obtained in Step 3.
      ```go
      const (
          spreadsheetID = "REPLACE_WITH_YOUR_SPREADSHEET_ID"
          sheetName     = "REPLACE_WITH_YOUR_SHEET_NAME"
          sheetId       = REPLACE_WITH_YOUR_GID
      )
      ```
3.  **Commit and Push** your changes.

Your attendance data will be automatically added to your Google Sheet, complete with table formatting.

### Customizing the Input Template
If you want to change the data columns that are sent to the Google Sheet, you need to edit two files:

1.  **`.github/workflows/write_to_sheet.yml`**:
    - Edit the `inputs:` section to add, remove, or change the fields that appear in the GitHub Actions form.
    - Make sure to also update the `env:` section to pass the new inputs to the Go script.

2.  **`main.go`**:
    - In the `main()` function, update the `os.Getenv` calls to read the new environment variables you set in the YAML file.
    - Adjust the `newRow` slice to match the new column structure of your Google Sheet. The order of variables in this slice must match the column order in your sheet.
      ```go
      // Example: Adjust this line in main.go
      newRow := sheets.ValueRange{
          Values: [][]interface{}{
              {nextNumberRow, todayDate, startTime, endTime, attendRecord, additionalInfo, newCustomField},
          },
      }
      ```
