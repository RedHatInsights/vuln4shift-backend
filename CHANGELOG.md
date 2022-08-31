# Changelog

<!--next-version-placeholder-->

## v0.13.0 (2022-08-31)
### Feature
* **manager:** Add sorting by cve severity for clusters ([`9c2c241`](https://github.com/RedHatInsights/vuln4shift-backend/commit/9c2c2410f23ea2932244db4d0d09f530df966fdf))

## v0.12.0 (2022-08-31)
### Feature
* **manager:** Add provider cluster filter ([`ac01502`](https://github.com/RedHatInsights/vuln4shift-backend/commit/ac015028edf639f78dc5bbf4ee3770812156f221))

## v0.11.0 (2022-08-31)
### Feature
* **manager:** Add version filter for clusters ([`0c4f39b`](https://github.com/RedHatInsights/vuln4shift-backend/commit/0c4f39b745d5f30266e3edebc6f258757afeda3e))

## v0.10.0 (2022-08-31)
### Feature
* **manager:** Add cluster status filter ([`6b05b4b`](https://github.com/RedHatInsights/vuln4shift-backend/commit/6b05b4b9777ccc61d26fa5d4576a4059ffee6168))

## v0.9.1 (2022-08-31)
### Fix
* **pyxis:** Dont init array with empty values ([`f129939`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f129939fe45639a3da3b8ee1f5c110df00c3a530))

## v0.9.0 (2022-08-30)
### Feature
* **manager:** Add report query for returning all records ([`4611008`](https://github.com/RedHatInsights/vuln4shift-backend/commit/46110081fd014c670c2c3a4bb6f49d6f1d2a25d3))

## v0.8.6 (2022-08-30)
### Fix
* **manager:** Cluster severity filter should be OR ([`9292213`](https://github.com/RedHatInsights/vuln4shift-backend/commit/9292213ad43c61498404235689e83b60db79bdfe))

## v0.8.5 (2022-08-29)
### Fix
* **digest-writer:** Collect digests from containers ([`536ebfc`](https://github.com/RedHatInsights/vuln4shift-backend/commit/536ebfc9ce846e6f0097debfaedea3d97a0ce88a))

## v0.8.4 (2022-08-26)
### Fix
* **manager:** Use provider+region str in list of all displayed providers ([`df18c7a`](https://github.com/RedHatInsights/vuln4shift-backend/commit/df18c7ad1861aa6898085b2f12f2a9575d11e60d))
* **manager:** Fetch metrics from AMS (adding metrics attribute doesnt't work) ([`96e90a3`](https://github.com/RedHatInsights/vuln4shift-backend/commit/96e90a3300568a428ce1e9738a0c4c4c0316fbb9))

## v0.8.3 (2022-08-26)
### Fix
* **pyxis:** Return unique image cves despite errata from pyxis ([`2ac645a`](https://github.com/RedHatInsights/vuln4shift-backend/commit/2ac645a16a178e5bb49f30f5d4f72c35d03a2918))

## v0.8.2 (2022-08-26)
### Fix
* **manager:** Revert this join, should be LEFT to display clusters without CVEs ([`f8928c0`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f8928c0130005054bca3dec8ebca7cceed67db26))

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
