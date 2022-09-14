-- ----------------------------------------------------------------------------
-- VULN4OPENSHIFT DB - Development Data Insert
--
-- This file inserts some data for initial development to use.
-- This is intended to be temporary until we get full flow working with
-- real data.
-- ----------------------------------------------------------------------------
INSERT INTO account (id, org_id) VALUES
(13, '013'),
(14, '014'),
(30, '030'),
(31, '031');

INSERT INTO cluster (id, uuid, status, version, provider, account_id, cve_cache_critical, cve_cache_important, cve_cache_moderate, cve_cache_low, workload) VALUES
(15, 'daac83ee-a390-420d-b892-cb9e1d006eca', 'ready'   , 'v_01', 'prov_1', 13, 0, 0, 0, 0, '{}'),
(16, '91eadf0a-3433-4862-8fc6-85e95ef7821c', 'sleeping', 'v_02', 'prov_2', 14, 0, 0, 0, 0, '{}'),
(17, '938a2b14-e6c0-4531-b203-377f9321afa6', 'up'      , 'v_03', 'prov_3', 14, 0, 0, 0, 0, '{}'),
(18, '83edaace-a390-420d-b892-cb9e1d006eca', 'ready'   , 'v_04', 'prov_4', 30, 0, 0, 0, 0, '{}'),
(19, 'adf91e0a-3433-4862-8fc6-85e95ef7821c', 'sleeping', 'v_05', 'prov_5', 31, 0, 0, 0, 0, '{}');

INSERT INTO arch (id, name) VALUES
(1, 'amd64');

INSERT INTO image (id, pyxis_id, modified_date, digest, arch_id) VALUES
(1, '57ea8d0d9c624c035f96f45e', '2022-04-01 11:29:57.343305+00', 'sha256:1643d7dd5ad3dd2a065221684cfd33eef84a5bc6025e080c05cde1755e1c15ea5c9d79561f92ce73517d6fd6d579e958543f9c3a61115b45863bbf7d42b2ac16454896d47f34287787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abb48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de45558428d0d193c05d941887cab4a62790e9d5d5f966f8f08ca05b736014181b5776386526629eb3fe706', 1),
(2, '57ea8d0e9c624c035f96f460', '2022-04-01 03:05:22.343305+00','sha256:58428d0d193c05ca05b736014181b5776386526629eb3f1643d7dd5ad3dd2a065221684cfd33eef84a5bc6025e080c05cde1755e1c15ea5c9d79561f92ce73517d6fd6d579e958543f9c3a61115b45863bbf7d42b2ac16454896d47f34287787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abd941887cab4a62790e9d5d5f966f8f08b48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de455e706', 1),
(3, '57ea8d0f9c624c035f96f466', '2022-03-27 00:00:57.343305+00','sha256:787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c91643d7dd5ad3dd2a065221684cfd33eef84a5bc6025e080c05cde1755e1c15ea5c9d79561f92ce73517d6fd6d579e958543f9c3a61115b45863bbf7d42b2ac16454896d47f342872f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abb48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de45558428d0d193c05d941887cab4a62790e9d5d5f966f8f08ca05b736014181b5776386526629eb3fe706', 1),
(4, '57ea8d109c624c035f96f468', '2022-03-30 14:14:57.343305+00','sha256:5bc6025e080c05cde1755e1c15ea5c9d79561f92ce73517d6fd6d579e9585431643d7dd5ad3dd2a065221684cfd33eef84af9c3a61115b45863bbf7d42b2ac16454896d47f34287787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abd941887cab4a62790e9d5d5f966f8f08b48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de45558428d0d193c05ca05b736014181b5776386526629eb3fe706', 1),
(5, '5f1eb63dbed8bd4f99e2a280', '2022-02-02 18:29:18.343305+00','sha256:87787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae81643d7dd5ad3dd2a065221684cfd33eef84a5bc6025e080c05cde1755e1c15ea5c9d79561f92ce73517d6fd6d579e958543f9c3a61115b45863bbf7d42b2ac16454896d47f34292e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abb48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de45558428d0d193c05d941887cab4a62790e9d5d5f966f8f08ca05b736014181b5776386526629eb3fe706', 1),
(6, '57ea8d149c624c035f96f47c', '2022-03-23 13:37:00.343305+00','sha256:19d79561f92ce73517d6fd6d579e958543f9c3a61115b45863bbf7d42b2ac16454896d47f34287787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e643d7dd5ad3dd2a065221684cfd33eef84a5bc6025e080c05cde1755e1c15ea5c7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abd941887cab4a62790e9d5d5f966f8f08b48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de45558428d0d193c05ca05b736014181b5776386526629eb3fe706', 1),
(7, '57ea8d159c624c035f96f47e', '2022-01-01 11:11:11.343305+00','sha256:cfd33eef84a5bc6025e080c05cde1755e1c15ea5c9d79561f92ce73517d6fd6d571643d7dd5ad3dd2a0652216849e958543f9c3a61115b45863bbf7d42b2ac16454896d47f34287787423aebca84dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abb48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de45558428d0d193c05d941887cab4a62790e9d5d5f966f8f08ca05b736014181b5776386526629eb3fe706', 1),
(8, '57ea8d389c624c035f96f524', '2022-02-22 22:22:22.343305+00','sha256:de1755e1c15ea5c9d79561f92ce73517d6fd6d579e958543f9c3a61115b45863bbf7d42b2ac16454896d47f34287787423aebca858428d0d193c05ca05b736014181b5776386526629eb3f1643d7dd5ad3dd2a065221684cfd33eef84a5bc6025e080c05c4dfd63a1cf303322901572666d58c51c9aa9070b94f360974bb506f2e7dfd3f80b7c44ac07436f13c1efaa0fd4258cd42ca037b2379e60a8944ae892e4d6e6a237defaade1c92f787d692c50be6491d71d4833373d3ce1ebd9eb6adf812db0d2f06709ecfd22ff218864ec9137493eb0310cdbaf8d45a134236b397316eff84c905994c76197609d25e50bf31934978797a9d3e2f85c65be8cf757d271ffca87a23da29596b3b5abd941887cab4a62790e9d5d5f966f8f08b48102e2c66790a1523442c4e479717d196125b90338c242e826ade289a859de455e706', 1),
(9, '57ea8d609c624c035f96f66f', '2022-03-05 05:14:23.343305+00','sha256:437939d5ff36ff4f7fe12bf1181059c3b41ff700dbc73c47321d57a168665fcbf2bac43dbf801ee24a1b7947d244fba10d8dfac0c49c77a004688e3a4e2093e036244e1334ac58fbf7576bcd24042d3130f533c6cdfb174191f4e3bcbf5a2f747ac832a908712a39549a30c0a0262b8a2f92cf26970bf0135144c24cf0f16b8351f81f168530c9807428e18eb4f74cff1c8484eb61d5e9b8e91fba95de783903b8fb3963402e9b49c6e71ed80f94d4522ab81798c59fb8fd03a7d6f968c9b82490569cff93e40616853c7d55e34d3e0c97fe60158b35b24fb469a11148b6cb638df90def6eee5e08ffdae6c2ace4c8c45bd2a0623a26979ffab56b812b5f4354f1d0e8397c1c1a88155833d7f5f372362a61ccd90d62f46c63a44730e7ae173b3c695b945fdb4d4b4b4700db3362e6d3d6d662af2870ab9698406eb3d0524264', 1),
(10, '57ea8d389c624c035f96f526', '2022-03-01 07:00:22.343305+00','sha256:47d1bcf86b085aacf96ad606aa44f6a9f67073a23643e7865a1e70761132cb8d978810363eb7da1580a59c5f0c0af767ba9431b84bad58f6843edc59e0561eba3685ce0f3feb95d022390d69d4eddccee902d276c1b3cf0c048fcae912485d83c6f85335c9f497862b0722f660966306b559f9ef9ce6d81d0c6c568fb09d04baef3d42e1677c55f5532d8c2b022ce1ac028519c65dc39badfa53fe035329f4ca61a3c6578a1b4a435d66568eef33d3a0af6edc2f3ee0fce69f9765c872fae0f007e67ce6875620bf6390e1c58b8ebc463eeb8272f403adeb4682652045e8409b8f94b2793e0842e73d390db93e70fee05109fde6ec134062755f055033a9cb3a9c941a06778fd6b6632e0be68c80bfd568c8a5d52fa1e4b78c45eb6a1db85d2283bc5748fc3079a43cc02c10a91122c31dd605e4e52303c710caa38578210199', 1),
(11, '57ea8d399c624c035f96f528', '2022-02-27 21:59:58.343305+00','sha256:410a55a2c85c7325e78b37f4a73e60610e6ffa18b6b1a1327a5b112ebaef3bd566dab2124b954d60c97dd84f21f53dd2457bb458d6ba3baa5814d00c875e908a6e1d8b60940f40282da254f147017590c01bf8fab28e63cef4cffaeb882d1c09aff761cdf4c506b13c7fe875fd06fa678f47c8f46d86c76fe295d67c19087f98696d9fe954f4afb1f1640dc85dbb8873f710f763de508973accc7578cd2ac0439aac38bed23a3ba5fd05c85774a25a497ae846a4721493971e63577617061c5f53a0f9220663aa2437993698dc875e02c00f3bc54d97a8783d240434b3c3750b17137a5aeaaeaff5be4f68a48d7bee002857955160e9742a6a4d0d5efb8d422db24cf38cbba6eaee01d13ba50086bc30768218ff0695fc09eac25dd6d5ada6b015ec73b3e1a50823f031a1e399c8d8db7b15d94ee748d7fba0a0ea34511b242a', 1),
(12, '57ea8d399c624c035f96f52a', '2022-01-26 23:59:00.343305+00','sha256:8f252831f96d1b3a9c516e084302d723d98a57f9185f6b3d5e58a95efe0f04acd2070278402ef085c6a63455065e8c1c3b696cc696271db2393884c2c9fdd1c0c4ea5ebac1c390a04b7fed2a41bf1d38f512739f1d4f5c19edc6c0a4fae2e8970ebd8d87756f8c4939e941e82904966c7b007d73fd99b2abb4a7573cd968178992294bea4b50394ef78abd536480437e189976d8cf34b47c9dbb53f1de6c571d4739a3891af63c91958572f74b9e98a3648c6f72ac247eacdd99800a4c111bd04d55dd34ffdfad1de19ba69a480c61839d8ff87988b80c54e2ee562128bc5ad6b702512b5db7415fa232130628630d57493a3d076a6863f31db4376c04d3830bbdf5e93254da222459f5e840057764bfdcd46b34822aecf138eaf04a47eb9ceef66b12aa1e905e951d202e59910b9c022c006fdba557a35e847d5c6bc0a91d77', 1);

INSERT INTO repository (id, pyxis_id, modified_date, registry, repository) VALUES
(1, '57ea8cd89c624c035f96f330', '2022-03-28 11:11:57.343305+00', 'registry.access.redhat.com', 'rhel7.1'),
(2, '57ea8cd79c624c035f96f328', '2022-02-21 20:57:11.343305+00', 'registry.access.redhat.com', 'rhel6.5'),
(3, '57ea8cd99c624c035f96f332', '2022-03-03 03:12:33.343305+00', 'registry.access.redhat.com', 'rhel7/sadc'),
(4, '57ea8cd99c624c035f96f334', '2022-02-02 02:33:44.343305+00', 'registry.access.redhat.com', 'rhel6');

INSERT INTO repository_image(repository_id, image_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 4),
(2, 5),
(2, 6),
(3, 7),
(3, 8),
(3, 9),
(4, 10),
(4, 11),
(4, 12);

INSERT INTO cluster_image (cluster_id, image_id) VALUES
(15, 1),
(16, 6),
(17, 6),
(18, 10),
(19, 11);

INSERT INTO cve (id, name, description, severity, cvss3_score, cvss3_metrics, cvss2_score, cvss2_metrics) VALUES
(20, 'CVE-2022-0001', 'mock CVE 01', 'High',      0.0, NULL,                                           8.0, 'AV:L/AC:H/Au:N/C:P/I:P/A:P'),
(21, 'CVE-2022-0002', 'mock CVE 02', 'Low',       1.5, 'CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:N/A:N', 0.0, NULL),
(22, 'CVE-2022-0003', 'mock CVE 03', 'Critical',  2.6, 'CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:N/A:N', 0.0, NULL),
(23, 'CVE-2022-0004', 'mock CVE 04', 'Moderate',  7.0, 'CVSS:3.0/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:N/A:N', 0.0, NULL),
(24, 'CVE-2022-0005', 'mock CVE 05', 'Important', 0.0, NULL,                                           6.0, 'AV:L/AC:H/Au:N/C:P/I:P/A:P');

INSERT INTO image_cve (image_id, cve_id) VALUES
(1, 20),
(2, 21),
(3, 22),
(3, 24),
(4, 20),
(4, 23),
(5, 23),
(6, 20),
(7, 20),
(8, 20),
(9, 22),
(10, 22),
(11, 23),
(11, 24);
