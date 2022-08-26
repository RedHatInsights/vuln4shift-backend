# Changelog

<!--next-version-placeholder-->

## v0.8.1 (2022-08-26)
### Fix
* **manager:** Make sure UUIDs from AMS API are valid ([`e8d6089`](https://github.com/RedHatInsights/vuln4shift-backend/commit/e8d608906abcad9dfe15c32360e94aa3c529e3ef))

## v0.8.0 (2022-08-25)
### Feature
* **manager:** Sync cluster details to db to allow sorting in DB ([`68fff0d`](https://github.com/RedHatInsights/vuln4shift-backend/commit/68fff0d888df6f0e5cc73d9149764dac16e7ec1f))
* **database:** Add missing columns to the cluster table and grant manager to update the table ([`a61ef5c`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a61ef5cd14693da2ee8319554523151ab0cd3ab7))

### Fix
* **digest-writer:** Missing cluster columns ([`7712b02`](https://github.com/RedHatInsights/vuln4shift-backend/commit/7712b02b8dc86e776135574ade4b6707e1d9ec7f))

## v0.7.0 (2022-08-24)
### Feature
* **manager:** Add unique sets of statuses, versions and providers ([`c2eaded`](https://github.com/RedHatInsights/vuln4shift-backend/commit/c2eadedaa082c6c07f2bbe1521db1abbfd22f916))

## v0.6.0 (2022-08-24)
### Feature
* **manager:** Add csv support for pageable endpoints ([`35c8526`](https://github.com/RedHatInsights/vuln4shift-backend/commit/35c8526770648c4295f308a7828fa6e178bcbb29))

## v0.5.0 (2022-08-19)
### Feature
* **manager:** Search clusters by display_name in AMS ([`eedd8c4`](https://github.com/RedHatInsights/vuln4shift-backend/commit/eedd8c4f30c3901db0bd7af45dd256f5d9054034))

## v0.4.0 (2022-08-19)
### Feature
* **manager:** AMS integration in CVE list ([`74becfb`](https://github.com/RedHatInsights/vuln4shift-backend/commit/74becfbf912a8516dc563c64c4e2978153e87215))
* **manager:** AMS integration in CVE exposed clusters ([`d3bbd88`](https://github.com/RedHatInsights/vuln4shift-backend/commit/d3bbd88139ed053f6b3679849407eb1090326f12))
* **manager:** AMS integration in cluster detail ([`81f96ac`](https://github.com/RedHatInsights/vuln4shift-backend/commit/81f96ac4e1e7ba4da0fadfbf87aa1f7515d5c3ce))
* **manager:** AMS integration in cluster list ([`9f86338`](https://github.com/RedHatInsights/vuln4shift-backend/commit/9f86338676e1a0687d98b2841cb4c53d1659b512))
* **manager:** Add AMS API client interface ([`7d47ca4`](https://github.com/RedHatInsights/vuln4shift-backend/commit/7d47ca4ceb3f25a9945ce32eb841db8052384e6d))

### Fix
* **manager:** Duplicate clusters when more than 1 image is affected by CVE ([`c89b05f`](https://github.com/RedHatInsights/vuln4shift-backend/commit/c89b05fe138e6f5fd5869215e50b9fd8a518c976))

## v0.3.0 (2022-08-16)
### Feature
* **manager-ams-api:** Prepare env variables ([`79a7591`](https://github.com/RedHatInsights/vuln4shift-backend/commit/79a75913c42a9eca1b42c577203a757e0372ae04))

## v0.2.2 (2022-08-15)
### Fix
* Make sure CVEs are ordered while inserting/deleting ([`13991c0`](https://github.com/RedHatInsights/vuln4shift-backend/commit/13991c0046eadbf9c7e63e46fd5282e2c16616e5))

## v0.2.1 (2022-08-05)
### Fix
* **manager:** Fix sorting for cvss2/3_score ([`8e91686`](https://github.com/RedHatInsights/vuln4shift-backend/commit/8e9168692778978c8a02c167002afe8966d1662c))

## v0.2.0 (2022-08-05)
### Feature
* **manager:** Add sort for display_name ([`fcd5c3e`](https://github.com/RedHatInsights/vuln4shift-backend/commit/fcd5c3e7a718b952f590786ba6300f5edd4b4960))

## v0.1.0 (2022-08-04)
### Feature
* Add cluster cleaner ([`99c6bf2`](https://github.com/RedHatInsights/vuln4shift-backend/commit/99c6bf270bc37d80508c93280bce2459a9819ef8))

## v0.0.1 (2022-07-29)
### Fix
* **manager:** Comma is allowed search character in fulltext search ([`f8dbd46`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f8dbd46e3914feb7a239cafba7c69252fbad67c1))
