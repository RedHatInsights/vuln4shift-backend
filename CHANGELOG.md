# Changelog

<!--next-version-placeholder-->

## v0.35.5 (2024-08-02)

### Fix

* **db:** Store registry_repository_version with index for searching ([`21f8b99`](https://github.com/RedHatInsights/vuln4shift-backend/commit/21f8b990511cb06dfc3fe21f3dc9e0a6d18e29f1))

## v0.35.4 (2024-07-31)

### Fix

* **gorm:** Force gorm to select specific fields ([`c619a4e`](https://github.com/RedHatInsights/vuln4shift-backend/commit/c619a4edbf16fae008fc644b3daa581794d85138))

## v0.35.3 (2024-06-06)

### Fix

* **expsync:** Possible SQL injection ([`52cad1e`](https://github.com/RedHatInsights/vuln4shift-backend/commit/52cad1e62a128e32185d4b3cbddf7b491ae309ee))

## v0.35.2 (2024-05-21)

### Fix

* **manager:** Base image count on correct distinct of columns ([`3227143`](https://github.com/RedHatInsights/vuln4shift-backend/commit/3227143887adc2cb50f1ad3aa8f72795dce19369))

## v0.35.1 (2024-05-15)

### Fix

* **digestwriter:** Refresh cluster CVE cache on each check-in ([#206](https://github.com/RedHatInsights/vuln4shift-backend/issues/206)) ([`3c74033`](https://github.com/RedHatInsights/vuln4shift-backend/commit/3c74033960b7b6270840e09e4805b4b8a421471f))

## v0.35.0 (2024-05-03)

### Feature

* **manager:** Registry filtering ([`eaeab77`](https://github.com/RedHatInsights/vuln4shift-backend/commit/eaeab774d150f2e2df664486df96433c48791e76))

## v0.34.11 (2024-05-03)

### Fix

* **manager:** Should be a column reference ([`a8de3af`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a8de3af87fc187bfe7ab6481c4c99d16de47ad3c))

## v0.34.10 (2024-05-03)

### Fix

* **manager:** Use version column + support sorting and filtering ([`f7bca47`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f7bca474effaa4eaa6792056989bf4dfaf5622a1))

## v0.34.9 (2024-05-03)

### Fix

* **pyxis:** Store the displayed version separately, this will make filtering and sorting a lot easier ([`1280bc8`](https://github.com/RedHatInsights/vuln4shift-backend/commit/1280bc86498137ba770ca23141148a28b8219685))

## v0.34.8 (2024-04-29)

### Fix

* Apply the numeric collation to repository/image names ([`2cddcbc`](https://github.com/RedHatInsights/vuln4shift-backend/commit/2cddcbc06b28a93964d9005c279c7af33100d33d))

## v0.34.7 (2024-04-29)

### Fix

* Add indexes to improve perf when selecting from CVE and image side ([`5455869`](https://github.com/RedHatInsights/vuln4shift-backend/commit/5455869ad1b3fc2ddfb577ffe00938f1c32758e3))

## v0.34.6 (2024-04-25)

### Fix

* **manager:** Count different aliases of an image ([`4fb30a8`](https://github.com/RedHatInsights/vuln4shift-backend/commit/4fb30a83c2a83c902d273340002328164634f844))

## v0.34.5 (2024-04-24)

### Fix

* **manager:** Count only clusters returned from amsclient ([`3ce03a2`](https://github.com/RedHatInsights/vuln4shift-backend/commit/3ce03a2cd9e431ee34c45ef71a5fa15f0add376f))

## v0.34.4 (2024-04-24)

### Fix

* **manager:** Fix and simplify cve images query ([`e604aa5`](https://github.com/RedHatInsights/vuln4shift-backend/commit/e604aa56fc7c73b75e3abfed931905534546947b))

## v0.34.3 (2024-04-24)

### Fix

* **manager:** Unique values of cluster images ([`76b0d15`](https://github.com/RedHatInsights/vuln4shift-backend/commit/76b0d1576a07f7ca6e2907a074e4876b6ad4b788))

## v0.34.2 (2024-04-23)

### Fix

* **/cves/[id]/exposed_images:** Fix query ([`64d7566`](https://github.com/RedHatInsights/vuln4shift-backend/commit/64d75661f8841f7f23189e2753808686f9c75bf2))

## v0.34.1 (2024-04-19)

### Fix

* **digestwriter:** Bump minimal assumed version, 0.10 doesn't seem to work anymore ([`b34cd43`](https://github.com/RedHatInsights/vuln4shift-backend/commit/b34cd43dc4ec659a6508c7ad9627cc6b5fd24852))

## v0.34.0 (2024-04-18)

### Feature

* Add /cves/[id]/exposed_images endpoint ([`586bb66`](https://github.com/RedHatInsights/vuln4shift-backend/commit/586bb6677ffe9fa22dad517f31e601e7e8d13486))

## v0.33.0 (2024-03-27)

### Feature

* **manager:** Add /clusters/[id]/exposed_images endpoint ([`c8e40c6`](https://github.com/RedHatInsights/vuln4shift-backend/commit/c8e40c608878c9a6ea50ab1fc837482e9e9fd688))

## v0.32.0 (2024-03-21)

### Feature

* Add tags to repository_image table ([`13cd3fb`](https://github.com/RedHatInsights/vuln4shift-backend/commit/13cd3fb9aed6d00d24c33526438f8b500e00ddbe))

## v0.31.0 (2024-02-27)

### Feature

* Add images_exposed to /cves ([`14b59b1`](https://github.com/RedHatInsights/vuln4shift-backend/commit/14b59b19a7574e8dff2add05e0c03d11ef681c9f))

## v0.30.0 (2024-02-27)

### Feature

* Use multiple kafka brokers ([`bacbad0`](https://github.com/RedHatInsights/vuln4shift-backend/commit/bacbad01c5ec37d5ad5d10890d9d540b78b72ac2))

## v0.29.2 (2024-02-26)

### Fix

* **database:** Sorting of cve names ([`f77c286`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f77c2864533e8a7e4a5555bb5ef6edeff70a59c6))

## v0.29.1 (2023-10-05)

### Fix

* **digestwriter:** Accept multiple account number types including empty strings ([`84734fc`](https://github.com/RedHatInsights/vuln4shift-backend/commit/84734fcec7b113cba9a9a0c80e37d22c34cf0e70))

## v0.29.0 (2023-09-19)

### Feature

* Ccx-sha-extractor decompression for messages ([#182](https://github.com/RedHatInsights/vuln4shift-backend/issues/182)) ([`188e9e4`](https://github.com/RedHatInsights/vuln4shift-backend/commit/188e9e4ab9b52af703add863f63c8f7a66a23c9f))

## v0.28.6 (2023-09-07)

### Fix

* Don't fetch Workload JSONB to optimize memory consumption ([`123ce24`](https://github.com/RedHatInsights/vuln4shift-backend/commit/123ce24ffa299f882ee233a57938f9cd475877ab))

## v0.28.5 (2023-09-04)

### Fix

* **digestwriter:** Accept different account number types ([`e5ae135`](https://github.com/RedHatInsights/vuln4shift-backend/commit/e5ae135db36f9e818640dcf94d078d4afab47196))

## v0.28.4 (2023-07-25)

### Fix

* **digestwriter:** Build cve cluster cache based on select from table ([`8093cf1`](https://github.com/RedHatInsights/vuln4shift-backend/commit/8093cf1ffd7086e0f7066f8096ddb233cc3bf011))

## v0.28.3 (2023-07-25)

### Fix

* Use older semantic-release ver ([`604931f`](https://github.com/RedHatInsights/vuln4shift-backend/commit/604931f4db7de14273dcc31bf4577bb8b8cd9a2f))

## v0.28.2 (2023-02-16)
### Fix
* **digestwriter:** Pass pt message by value ([`85406bb`](https://github.com/RedHatInsights/vuln4shift-backend/commit/85406bb6dcbbfddfb571d10655393999c7517b03))

## v0.28.1 (2023-02-16)
### Fix
* **base:** Turn on sarama success return values ([`a82743e`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a82743eba53ae0a4a15a2af78fc03940ec970b05))
* **digestwriter:** Send all payload messages in goroutine ([`1ea978a`](https://github.com/RedHatInsights/vuln4shift-backend/commit/1ea978a042a3bd032a751f0029b9a3bc19111844))

## v0.28.0 (2023-02-02)
### Feature
* **digestwriter:** Add Payload Tracker metrics ([`71b0f58`](https://github.com/RedHatInsights/vuln4shift-backend/commit/71b0f587dce5907fdb32f7278875fd8a3740775d))

## v0.27.0 (2023-01-30)
### Feature
* **digestwriter:** Add Payload Tracker feature flag ([`053bfb6`](https://github.com/RedHatInsights/vuln4shift-backend/commit/053bfb64a842b80cde9be74cc556c5ebe563b1f1))

## v0.26.2 (2023-01-30)
### Fix
* Change apiPath to ocp-vulnerability ([`a152278`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a152278bda32413c978d7934479aa8bda421e916))

## v0.26.1 (2023-01-26)
### Fix
* **manager:** Fix references to columns in removed subquery ([`f382277`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f3822774be86d68897c7bfe949795f8f0a0ea2f9))

## v0.26.0 (2023-01-16)
### Feature
* **digestwriter:** Add Payload Tracker kafka producer ([`4f039e3`](https://github.com/RedHatInsights/vuln4shift-backend/commit/4f039e3b9d009d37b9839a0a2bf43763c5131047))
* **base:** Add Payload Tracker event ([`18b4cc5`](https://github.com/RedHatInsights/vuln4shift-backend/commit/18b4cc54bba10a01a8b7721c210cec3b544f6080))
* **base:** Add Kafka producer ([`cdc203c`](https://github.com/RedHatInsights/vuln4shift-backend/commit/cdc203c0afb24a0ba8c8f07a3a26fbe78fe69434))

## v0.25.1 (2022-12-16)
### Fix
* **expsync:** Add missing clowdapp env vars ([`5efc954`](https://github.com/RedHatInsights/vuln4shift-backend/commit/5efc954d32c3bf449bc7f35e9aff5d74c25b6710))

## v0.25.0 (2022-12-16)
### Feature
* **tests:** Init caches on unit test db ([`92a6803`](https://github.com/RedHatInsights/vuln4shift-backend/commit/92a6803640969f1b2fc6f8765c7499bcddfd5a64))
* **manager:** Start returning cluster-cve cached counts ([`0d3fcf7`](https://github.com/RedHatInsights/vuln4shift-backend/commit/0d3fcf7e1490b378c49473a57bbae81ec15e3dca))

## v0.24.0 (2022-12-15)
### Feature
* Log to cloudwatch (+minor refactor) ([`cd71d29`](https://github.com/RedHatInsights/vuln4shift-backend/commit/cd71d2984b6f03109616d9bff045b0552bc0b599))

## v0.23.0 (2022-12-15)
### Feature
* **pyxis:** Validate cves from pyxis ([`a6023c3`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a6023c3a295872ce87a681b28a11dea80e3b1371))

## v0.22.0 (2022-12-12)
### Feature
* **manager:** Add exposed clusters count endpoint ([`387729a`](https://github.com/RedHatInsights/vuln4shift-backend/commit/387729a9d847626ce6b744cfe8a0251785c0f3a8))

## v0.21.0 (2022-12-12)
### Feature
* **manager:** Add exploits field to cluster cves ([`c9200d2`](https://github.com/RedHatInsights/vuln4shift-backend/commit/c9200d2e5e22c3639bbd244b3e14d75b2ec84f71))
* **manager:** Add exploits field to cve details endpoint ([`c4eb52b`](https://github.com/RedHatInsights/vuln4shift-backend/commit/c4eb52ba4231bcd961bec58b699992e9f374c2d0))
* **manager:** Add exploits field to cve endpoint ([`2aee3f9`](https://github.com/RedHatInsights/vuln4shift-backend/commit/2aee3f9511006627733f6b73d0c3fc24b7469b24))
* **expsync:** Add exploit sync job ([`cdc5202`](https://github.com/RedHatInsights/vuln4shift-backend/commit/cdc5202288aa108dbdf16a9b87db393b57b54ee5))

## v0.20.0 (2022-11-28)
### Feature
* **digestwriter:** Add cve-cluster cache calculation ([`495bb12`](https://github.com/RedHatInsights/vuln4shift-backend/commit/495bb1238ce2471729cdefab01eff756ccae2b67))

## v0.19.6 (2022-10-14)
### Fix
* **pyxis:** Improve repo images check ([`f073940`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f073940fc1a831672a1c4b0707d221a94dc887b6))

## v0.19.5 (2022-10-11)
### Fix
* **pyxis:** Correct registerMissingCves metric incrementation ([`35fe26d`](https://github.com/RedHatInsights/vuln4shift-backend/commit/35fe26d6c6d940c770596e38e8c1628c14bf7731))

## v0.19.4 (2022-10-03)
### Fix
* **manager:** Filter *_all lists for a single CVE scope ([`a9f13e2`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a9f13e233ee9508095f2bff05450f94bb3fc24d7))

## v0.19.3 (2022-09-26)
### Fix
* **digestwriter:** Unlink images from cluster ([`6baa5e3`](https://github.com/RedHatInsights/vuln4shift-backend/commit/6baa5e35f2e479770d349eaa722c6dd3090dbf61))

## v0.19.2 (2022-09-22)
### Fix
* **pyxis:** Use both pyxis_id and reg/repo key for sync ([`3157364`](https://github.com/RedHatInsights/vuln4shift-backend/commit/315736474f2a1d771d1470f4ae2bc6eee1e91474))

## v0.19.1 (2022-09-21)
### Fix
* **manager:** Set max header size to 65535 bytes ([`4fa6393`](https://github.com/RedHatInsights/vuln4shift-backend/commit/4fa639376e968adbc01589d08e7cf1ceced45257))

## v0.19.0 (2022-09-20)
### Feature
* **vmsync:** Prune CVEs during sync ([`79bb0e0`](https://github.com/RedHatInsights/vuln4shift-backend/commit/79bb0e0a788a2c82e8784b90914e93008e270def))

## v0.18.0 (2022-09-20)
### Feature
* **pyxis:** Skip image CVE sync (debug purposes) ([`b999034`](https://github.com/RedHatInsights/vuln4shift-backend/commit/b99903477394d2d2cc26bf43ed7787adc84a3eac))

## v0.17.3 (2022-09-19)
### Fix
* **pyxis:** Referencing attributes of apiImage will use only the last value in the loop, reference local variables instead ([`7b7bc7e`](https://github.com/RedHatInsights/vuln4shift-backend/commit/7b7bc7e8aa60ef5e02479391c161e31b4b671c26))

## v0.17.2 (2022-09-19)
### Fix
* Check more image digests ([`8dbb70f`](https://github.com/RedHatInsights/vuln4shift-backend/commit/8dbb70f3bd15d93f7b8badb3729ca42f100e9da6))

## v0.17.1 (2022-09-15)
### Fix
* **manager:** Do not encode undefined value ([`e7489f4`](https://github.com/RedHatInsights/vuln4shift-backend/commit/e7489f468479dd47d6258ee7d8a6910e908c5822))

## v0.17.0 (2022-09-15)
### Feature
* **digestwriter:** Store cluster workload JSON ([`73560ee`](https://github.com/RedHatInsights/vuln4shift-backend/commit/73560ee652d19c6a4ee6e41e0d31d55a99b9b094))

## v0.16.1 (2022-09-08)
### Fix
* **pyxis:** Use repository map with reg,repo key ([`2f2c268`](https://github.com/RedHatInsights/vuln4shift-backend/commit/2f2c26887280453de96e2c7e1b4c49528080514f))

## v0.16.0 (2022-09-08)
### Feature
* **manager:** Support status, version and provider filters in exposed clusters API ([`f978366`](https://github.com/RedHatInsights/vuln4shift-backend/commit/f978366717909f24e6baffbb40d100c1b52c844f))

## v0.15.0 (2022-09-07)
### Feature
* Drop account_number ([`64b1608`](https://github.com/RedHatInsights/vuln4shift-backend/commit/64b160871d206f7eaf904cf6fe57e0841387382b))
* **digestwriter:** Drop usage of account number ([`cccae2f`](https://github.com/RedHatInsights/vuln4shift-backend/commit/cccae2f2060f44cbfc64ffa89e5e3029b306b3ad))

## v0.14.7 (2022-09-06)
### Fix
* **manager:** Sort versions by number ([`b6cb581`](https://github.com/RedHatInsights/vuln4shift-backend/commit/b6cb5811d2aae4042d1ec86ff573c0838daad3f2))

## v0.14.6 (2022-09-06)
### Fix
* **kafka:** Create TLS context despite empty cacert ([`a55bd38`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a55bd38bf153f425832ff63f6017db417eca936d))

## v0.14.5 (2022-09-06)
### Fix
* **database:** Sorting of cluster versions ([`885e5c4`](https://github.com/RedHatInsights/vuln4shift-backend/commit/885e5c4423c9b65ec7557e708de7332af65cb17c))

## v0.14.4 (2022-09-06)
### Fix
* **manager:** Report query should use parsedValues ([`77ff00a`](https://github.com/RedHatInsights/vuln4shift-backend/commit/77ff00aa1991f0e4d00784f1700af1e810cf5d31))

## v0.14.3 (2022-09-05)
### Fix
* **manager:** Return meta lists in predictable order ([`feb0cdc`](https://github.com/RedHatInsights/vuln4shift-backend/commit/feb0cdcadc910e8b5bcc8594e5a0b279c45914b5))

## v0.14.2 (2022-09-05)
### Fix
* **manager:** Support comma-separated values in status filter ([`e2025cc`](https://github.com/RedHatInsights/vuln4shift-backend/commit/e2025cc333c2617a22100b31ed6340e6c413b025))

## v0.14.1 (2022-09-03)
### Fix
* **pyxis:** Use docker_image_digest if manifest_list_digest is empty ([`4666f1d`](https://github.com/RedHatInsights/vuln4shift-backend/commit/4666f1d3ac258bbe283da1b51877349cdf00b466))

## v0.14.0 (2022-09-02)
### Feature
* **pyxis:** Store image archs ([`a87d462`](https://github.com/RedHatInsights/vuln4shift-backend/commit/a87d46201d54bd3b5aa199a5f1c3a90ffbb3cede))
* **database:** Add arch for image and remove digest unique constaint -> manifest list digest will be used ([`db4442d`](https://github.com/RedHatInsights/vuln4shift-backend/commit/db4442df48cb40ea5662f3762a38930a85b44109))
* **pyxis:** Support forcing the sync ([`1b10000`](https://github.com/RedHatInsights/vuln4shift-backend/commit/1b10000e5b114f7c784384bf36da913131628c5a))

### Fix
* **digest-writer:** Match images with given arch ([`1db8d19`](https://github.com/RedHatInsights/vuln4shift-backend/commit/1db8d192d434d5db4274e58db1479a91eca799a4))
* **pyxis:** Store manifest_list_digest instead ([`dc0ae6b`](https://github.com/RedHatInsights/vuln4shift-backend/commit/dc0ae6b02cfdeef43cf57b1d09e5e81a6b749c7b))

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
