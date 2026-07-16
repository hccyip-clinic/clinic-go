# Hong Ching Clinic (HCC) - Data Model Requirements

## 1. Entity: Official Receipt (正式收據)

**Purpose**: Store official receipt / payment records for treatments.

### Core Fields

| Field Name (English)        | Chinese Name        | Type               | Required | Description                                          |
| --------------------------- | ------------------- | ------------------ | -------- | ---------------------------------------------------- |
| receipt_number              | 收據編號            | String             | Yes      | Receipt number (e.g. "RCPT-20260708143022-ABCDE" or "0000206")                      |
| issue_date                  | 日期                | Date               | Yes      | Date when the receipt was issued                     |
| patient_name                | 患者姓名 / 先生女士 | String             | Yes      | Full name of the patient                             |
| patient_gender              | 性別                | Enum (Male/Female) | No       | Patient gender                                       |
| practitioner_name           | 中醫師姓名          | String             | Yes      | Name of the registered Chinese medicine practitioner |
| practitioner_reg_no         | 註冊編號            | String             | Yes      | Registered practitioner number (e.g. "002582")       |
| treatment_tui_na            | 經跌打 / 推拿       | Boolean            | No       | Whether Osteopathic Tui-Na treatment was provided    |
| treatment_acupuncture       | 針灸                | Boolean            | No       | Whether Acupuncture treatment was provided           |
| treatment_internal_medicine | 內科治療            | Boolean            | No       | Whether Internal Medicine treatment was provided     |
| total_amount_hkd            | 合共 HKD            | Decimal(10,2)      | Yes      | Total amount charged in Hong Kong Dollars            |
| received_by                 | 收款人              | String             | No       | Name of the staff who received the payment           |
| clinic_stamp_applied        | 診所印章            | Boolean            | No       | Indicates if the official clinic stamp was applied   |
| clinic_address              | 地址                | String             | Yes      | Full clinic address                                  |
| clinic_telephone            | 電話                | String             | Yes      | Clinic telephone number                              |

---

## 2. Entity: Sick Leave / Medical Certificate (病假 / 到診證明書)

**Purpose**: Record sick leave recommendations and consultation certificates.

### Core Fields

| Field Name (English)   | Chinese Name | Type               | Required | Description                                       |
| ---------------------- | ------------ | ------------------ | -------- | ------------------------------------------------- |
| certificate_number     | 證明書編號   | String             | No       | Medical certificate serial number (if applicable) |
| patient_name           | 姓名         | String             | Yes      | Full name of the patient                          |
| patient_gender         | 性別         | Enum (Male/Female) | No       | Patient gender                                    |
| patient_age            | 年齡         | Integer            | No       | Patient age in years                              |
| patient_hkid           | 身份證 No.   | String             | No       | Patient HKID number                               |
| diagnosis              | 診斷 / 主証  | String / Text      | Yes      | Diagnosis or main clinical evidence               |
| is_sick_leave          | 病假         | Boolean            | No       | Whether sick leave is recommended                 |
| sick_leave_days        | 建議休息天數 | Integer            | No       | Number of recommended rest days                   |
| sick_leave_from_date   | 由...日起    | Date               | No       | Start date of the recommended sick leave          |
| sick_leave_to_date     | 至...日止    | Date               | No       | End date of the recommended sick leave            |
| consultation_datetime  | 到診時間     | DateTime           | No       | Date and time of the consultation / visit         |
| remarks                | 備註         | Text               | No       | Additional remarks or notes                       |
| practitioner_signature | 中醫師簽名   | String / ImageRef  | Yes      | Signature of the Chinese medicine practitioner    |
| practitioner_name      | 中醫師姓名   | String             | Yes      | Name of the registered practitioner               |
| practitioner_reg_no    | 註冊編號     | String             | Yes      | Registered practitioner number (e.g. "002582")    |
| issue_date             | 簽發日期     | Date               | Yes      | Date when the medical certificate was issued      |
| clinic_address         | 地址         | String             | Yes      | Full clinic address                               |
| clinic_telephone       | 電話         | String             | Yes      | Clinic telephone number                           |

---

## 3. Shared Clinic Information

| Field Name (English) | Chinese Name | Type   | Description                               |
| -------------------- | ------------ | ------ | ----------------------------------------- |
| clinic_name_en       | 診所英文名稱 | String | Clinic English name                       |
| clinic_name_zh       | 診所中文名稱 | String | Clinic Chinese name                       |
| clinic_reg_no        | 註冊編號     | String | Clinic / Practitioner registration number |
| clinic_full_address  | 地址         | String | Full clinic address (bilingual)           |
| clinic_telephone     | 電話         | String | Clinic contact telephone number           |

---

## 4. Treatment Types (Lookup)

- `treatment_tui_na` (經跌打 / 推拿)
- `treatment_acupuncture` (針灸)
- `treatment_internal_medicine` (內科治療)

---

## 5. Recommended Relationships

- One `Patient` → Many `Receipts`
- One `Patient` → Many `MedicalCertificates`
- One `Practitioner` → Many `Receipts` and `MedicalCertificates`

---

## 6. Additional Recommendations

- Use UUID as primary keys for all main entities.
- Store image references for practitioner signatures and clinic stamps.
- Support bilingual (English + Traditional Chinese) display and output.
- Include audit fields: `created_at`, `updated_at`, `created_by`.

This version uses **English field names** as the primary column, with **Chinese names** and full **English descriptions**.

---

## 7. Receipt Number Format

### New Format (Timestamp-based)
**Format:** `RCPT-yyyyMMddHHmmss-ABCDE`
- **Prefix:** `RCPT-` (fixed)
- **Timestamp:** 14-digit local time (yyyyMMddHHmmss)
- **Suffix:** 5 uppercase random letters (A-Z)
- **Example:** `RCPT-20260708143022-ABCDE`
- **Total Length:** 24 characters

### Legacy Format (Sequential)
**Format:** 7-digit zero-padded integer
- **Example:** `0000206`
- **Status:** Still valid for existing receipts

Both formats are supported by `isValidReceiptNumber()` validation.