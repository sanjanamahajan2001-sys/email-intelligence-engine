PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		base_email TEXT,
		has_alias BOOLEAN,
		is_valid BOOLEAN,
		syntax BOOLEAN,
		dns BOOLEAN,
		smtp BOOLEAN,
		disposable BOOLEAN,
		role BOOLEAN,
		domain_age_years REAL,
		reputation_score INTEGER,
		risk_level TEXT,
		provider TEXT,
		tld_trust TEXT,
		source TEXT,
		catch_all BOOLEAN,
		greylisted BOOLEAN,
		identity_age_years INTEGER,
		confidence_score INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	, engagement_probability INTEGER DEFAULT 0, last_smtp_response TEXT, engagement_insight TEXT, lifecycle_state TEXT, engagement_factors TEXT, message TEXT);
INSERT INTO scans VALUES(559,'active-stale@example.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,95,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 05:37:31',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(560,'catchall-stale@example.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,80,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-14 05:37:31',0,NULL,NULL,'CATCH-ALL',NULL,NULL);
INSERT INTO scans VALUES(561,'dispo-fresh@example.com',NULL,NULL,0,NULL,NULL,NULL,NULL,NULL,NULL,10,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-29 05:37:31',0,NULL,NULL,'DISPOSABLE',NULL,NULL);
INSERT INTO scans VALUES(562,'active-stale@example.com','active-stale@example.com',0,0,1,1,0,0,0,30.6576775953303801,35,'High','','High','CLI-UPDATE',0,0,30.6576775953303801,60,'2026-04-03 05:37:37',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(563,'catchall-stale@example.com','catchall-stale@example.com',0,0,1,1,0,0,0,30.6576775953303801,35,'High','','High','CLI-UPDATE',0,0,30.6576775953303801,60,'2026-04-03 05:37:37',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(564,'contact@stripe.com','contact@stripe.com',0,1,1,1,1,0,1,30.5763264203471365,100,'Low','Tier-1 Enterprise','High','CLI',0,0,30.5763264203471365,60,'2026-04-03 05:38:11',60,'250 2.1.5 OK','Moderate likelihood. Common for enterprise or catch-all addresses.','ACTIVE','["+15: Verified Enterprise Provider","+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)","-15: Role-based address (info@/admin@)"]','Verified Corporate Identity');
INSERT INTO scans VALUES(565,'ceo@apple.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,95,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(566,'support@stripe.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,80,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-14 12:20:08',0,NULL,NULL,'CATCH-ALL',NULL,NULL);
INSERT INTO scans VALUES(567,'marketing@google.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,90,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-26 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(568,'noreply@microsoft.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,90,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-19 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(569,'hello@github.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,85,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-09 12:20:08',0,NULL,NULL,'CATCH-ALL',NULL,NULL);
INSERT INTO scans VALUES(570,'sales@amazon.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,95,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(571,'dev@meta.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,95,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-22 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(572,'jobs@facebook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,95,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-25 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(573,'info@netflix.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,80,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-12 12:20:08',0,NULL,NULL,'CATCH-ALL',NULL,NULL);
INSERT INTO scans VALUES(574,'legal@twitter.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,90,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-04 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(575,'loadtest-1@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(576,'loadtest-2@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(577,'loadtest-3@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(578,'loadtest-4@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(579,'loadtest-5@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(580,'loadtest-6@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(581,'loadtest-7@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(582,'loadtest-8@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(583,'loadtest-9@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(584,'loadtest-10@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(585,'loadtest-11@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(586,'loadtest-12@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(587,'loadtest-13@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(588,'loadtest-14@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(589,'loadtest-15@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(590,'loadtest-16@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(591,'loadtest-17@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(592,'loadtest-18@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(593,'loadtest-19@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(594,'loadtest-20@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(595,'loadtest-21@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(596,'loadtest-22@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(597,'loadtest-23@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(598,'loadtest-24@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(599,'loadtest-25@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(600,'loadtest-26@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(601,'loadtest-27@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(602,'loadtest-28@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(603,'loadtest-29@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(604,'loadtest-30@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(605,'loadtest-31@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(606,'loadtest-32@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(607,'loadtest-33@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(608,'loadtest-34@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(609,'loadtest-35@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(610,'loadtest-36@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(611,'loadtest-37@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(612,'loadtest-38@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(613,'loadtest-39@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(614,'loadtest-40@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(615,'loadtest-41@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(616,'loadtest-42@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(617,'loadtest-43@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(618,'loadtest-44@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(619,'loadtest-45@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(620,'loadtest-46@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(621,'loadtest-47@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(622,'loadtest-48@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(623,'loadtest-49@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(624,'loadtest-50@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(625,'loadtest-51@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(626,'loadtest-52@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(627,'loadtest-53@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(628,'loadtest-54@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(629,'loadtest-55@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(630,'loadtest-56@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(631,'loadtest-57@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(632,'loadtest-58@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(633,'loadtest-59@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(634,'loadtest-60@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(635,'loadtest-61@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(636,'loadtest-62@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(637,'loadtest-63@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(638,'loadtest-64@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(639,'loadtest-65@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(640,'loadtest-66@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(641,'loadtest-67@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(642,'loadtest-68@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(643,'loadtest-69@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(644,'loadtest-70@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(645,'loadtest-71@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(646,'loadtest-72@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(647,'loadtest-73@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(648,'loadtest-74@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(649,'loadtest-75@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(650,'loadtest-76@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(651,'loadtest-77@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(652,'loadtest-78@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(653,'loadtest-79@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(654,'loadtest-80@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(655,'loadtest-81@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(656,'loadtest-82@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(657,'loadtest-83@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(658,'loadtest-84@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(659,'loadtest-85@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(660,'loadtest-86@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(661,'loadtest-87@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(662,'loadtest-88@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(663,'loadtest-89@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(664,'loadtest-90@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(665,'loadtest-91@outlook.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(666,'loadtest-92@yahoo.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(667,'loadtest-93@protonmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(668,'loadtest-94@zoho.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(669,'loadtest-95@icloud.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(670,'loadtest-96@gmx.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(671,'loadtest-97@mail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(672,'loadtest-98@fastmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(673,'loadtest-99@yandex.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(674,'loadtest-100@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-24 12:20:08',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(675,'loadtest-26@gmx.com','loadtest-26@gmx.com',0,0,1,1,0,0,0,31.9289475005568306,35,'High','','High','CLI-UPDATE',0,0,31.9289475005568306,60,'2026-04-03 12:20:26',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(676,'loadtest-16@gmx.com','loadtest-16@gmx.com',0,0,1,1,0,0,0,31.9289475005568306,35,'High','','High','CLI-UPDATE',0,0,31.9289475005568306,60,'2026-04-03 12:20:26',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(677,'loadtest-17@mail.com','loadtest-17@mail.com',0,0,1,1,0,0,0,29.046641539926135,35,'High','','High','CLI-UPDATE',0,0,29.046641539926135,60,'2026-04-03 12:20:26',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(678,'loadtest-27@mail.com','loadtest-27@mail.com',0,0,1,1,0,0,0,29.046641539926135,35,'High','','High','CLI-UPDATE',0,0,29.046641539926135,60,'2026-04-03 12:20:26',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(679,'loadtest-19@yandex.com','loadtest-19@yandex.com',0,0,1,1,0,0,0,27.5426461316319262,35,'High','','High','CLI-UPDATE',0,0,27.5426461316319262,60,'2026-04-03 12:20:27',0,'550 5.7.1 No such user! 1775218826-PKcilu1KS4Y0-BLWnbI3V','This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(680,'loadtest-29@yandex.com','loadtest-29@yandex.com',0,0,1,1,0,0,0,27.5426461316319262,35,'High','','High','CLI-UPDATE',0,0,27.5426461316319262,60,'2026-04-03 12:20:28',0,'550 5.7.1 No such user! 1775218826-PKcSsu1LEmI0-CorVjV09','This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(681,'loadtest-100@gmail.com','loadtest-100@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI-UPDATE',0,0,30.6584990413732718,60,'2026-04-03 12:20:28',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser d9443c01a7336-2b274b3f9f0si145040935ad.168 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(682,'loadtest-10@gmail.com','loadtest-10@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI-UPDATE',0,0,30.6584990413732718,60,'2026-04-03 12:20:28',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser d2e1a72fcca58-82cf9cc4e21si14930890b3a.97 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(683,'loadtest-30@gmail.com','loadtest-30@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI-UPDATE',0,0,30.6584990413732718,60,'2026-04-03 12:20:28',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser 41be03b00d2f7-c76c648cc8asi11079287a12.81 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(684,'info@netflix.com','info@netflix.com',0,1,1,1,1,0,1,28.4110176760871908,100,'Low','','High','CLI-UPDATE',0,0,28.4110176760871908,60,'2026-04-03 12:20:28',60,'250 2.1.5 OK','Moderate likelihood. Common for enterprise or catch-all addresses.','ACTIVE','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)","+10: Consistent Historical Activity","-15: Role-based address (info@/admin@)"]','Established Legacy Identity');
INSERT INTO scans VALUES(685,'loadtest-25@icloud.com','loadtest-25@icloud.com',0,0,1,1,0,0,0,27.2329429458297482,35,'High','Global Consumer','High','CLI-UPDATE',0,0,27.2329429458297482,60,'2026-04-03 12:20:28',35,'550 5.7.1 Mail from IP 223.233.81.126 was rejected due to listing in Spamhaus PBL. For details please see http://www.spamhaus.org/query/bl?ip=223.233.81.126','Low likelihood of reply due to infrastructure or identity age signals.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","-30: SMTP Delivery Blocked (Security Filtering)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(686,'loadtest-20@gmail.com','loadtest-20@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI-UPDATE',0,0,30.6584990413732718,60,'2026-04-03 12:20:28',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser d9443c01a7336-2b2749fdb7esi118489165ad.147 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(687,'legal@twitter.com','legal@twitter.com',0,1,1,1,1,0,0,26.2151950631222341,100,'Low','','High','CLI-UPDATE',0,0,26.2151950631222341,60,'2026-04-03 12:20:28',75,'250 2.1.5 OK','High likelihood of reply from an established infrastructure.','ACTIVE','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)","+10: Consistent Historical Activity"]','Established Legacy Identity');
INSERT INTO scans VALUES(688,'hello@github.com','hello@github.com',0,0,1,1,0,0,0,18.4944255706368103,35,'High','Tier-1 Enterprise','High','CLI-UPDATE',0,0,18.4944255706368103,60,'2026-04-03 12:20:28',30,'550 5.7.1 Service unavailable, Client host [223.233.81.126] blocked using Spamhaus. To request removal from this list see https://www.spamhaus.org/query/ip/223.233.81.126 AS(1450) [SJ1PEPF00001CE1.namprd05.prod.outlook.com 2026-04-03T12:20:26.464Z 08DE8FAEF80F4FFA]','Low likelihood of reply due to infrastructure or identity age signals.','ABANDONED','["+15: Verified Enterprise Provider","+20: Established Legacy Identity (\u003e5 yrs)","-15: Temporary Bounce/Greylisting detected","-30: SMTP Delivery Blocked (Security Filtering)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(689,'loadtest-2@yahoo.com','loadtest-2@yahoo.com',0,0,1,1,0,0,0,31.2274560669390162,35,'High','Global Consumer','High','CLI-UPDATE',0,0,31.2274560669390162,60,'2026-04-03 12:20:29',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(690,'loadtest-15@icloud.com','loadtest-15@icloud.com',0,0,1,1,0,0,0,27.2329429458297482,35,'High','Global Consumer','High','CLI-UPDATE',0,0,27.2329429458297482,60,'2026-04-03 12:20:29',35,'550 5.7.1 Mail from IP 223.233.81.126 was rejected due to listing in Spamhaus PBL. For details please see http://www.spamhaus.org/query/bl?ip=223.233.81.126','Low likelihood of reply due to infrastructure or identity age signals.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","-30: SMTP Delivery Blocked (Security Filtering)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(691,'loadtest-22@yahoo.com','loadtest-22@yahoo.com',0,0,1,1,0,0,0,31.2274560669390162,35,'High','Global Consumer','High','CLI-UPDATE',0,0,31.2274560669390162,60,'2026-04-03 12:20:29',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(692,'loadtest-14@zoho.com','loadtest-14@zoho.com',0,0,1,1,0,0,0,22.2258479015073788,35,'High','Zoho Workplace','High','CLI-UPDATE',0,0,22.2258479015073788,60,'2026-04-03 12:20:29',0,'550 5.1.1 User does not exist - <loadtest-14@zoho.com>','This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(693,'loadtest-12@yahoo.com','loadtest-12@yahoo.com',0,0,1,1,0,0,0,31.2274560669390162,35,'High','Global Consumer','High','CLI-UPDATE',0,0,31.2274560669390162,60,'2026-04-03 12:20:29',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(694,'dev@meta.com','dev@meta.com',0,1,1,1,1,0,0,35.2219766763897439,80,'Medium','','High','CLI-UPDATE',1,0,35.2219766763897439,60,'2026-04-03 12:20:31',35,'250 2.1.5 OK','Low likelihood of reply due to infrastructure or identity age signals.','CATCH-ALL','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)","+10: Consistent Historical Activity","-40: Catch-all domain (Reduced Delivery Confidence)"]','Risky: Domain is a Catch-All Configuration');
INSERT INTO scans VALUES(695,'loadtest-28@fastmail.com','loadtest-28@fastmail.com',0,0,1,1,0,0,0,31.3370525941345796,35,'High','','High','CLI-UPDATE',0,0,31.3370525941345796,60,'2026-04-03 12:20:31',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(696,'ceo@apple.com','ceo@apple.com',0,1,1,1,1,0,0,39.145265854102881,100,'Low','Global Consumer','High','CLI-UPDATE',0,0,39.145265854102881,60,'2026-04-03 12:20:31',75,'250 2.1.5 OK','High likelihood of reply from an established infrastructure.','ACTIVE','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)","+10: Consistent Historical Activity"]','Trusted: Legacy Consumer Provider');
INSERT INTO scans VALUES(697,'loadtest-18@fastmail.com','loadtest-18@fastmail.com',0,0,1,1,0,0,0,31.3370525941345796,35,'High','','High','CLI-UPDATE',0,0,31.3370525941345796,60,'2026-04-03 12:20:32',65,'','Moderate likelihood. Common for enterprise or catch-all addresses.','ABANDONED','["+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(698,'jobs@facebook.com','jobs@facebook.com',0,0,1,1,0,0,0,29.0329430923180673,35,'High','','High','CLI-UPDATE',0,0,29.0329430923180673,60,'2026-04-03 12:20:32',0,'550 5.1.1 RCP-P1 Domain facebook.com no longer available https://www.facebook.com/postmaster/response_codes?ip=223.233.81.126#RCP-P1','This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(699,'loadtest-11@outlook.com','loadtest-11@outlook.com',0,0,1,1,0,0,0,31.6467482405604805,35,'High','Microsoft 365','High','CLI-UPDATE',0,0,31.6467482405604805,60,'2026-04-03 12:20:32',80,'','High likelihood of reply from an established infrastructure.','ABANDONED','["+25: Tier-1 Infrastructure (Google/Microsoft)","+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(700,'loadtest-31@outlook.com','loadtest-31@outlook.com',0,0,1,1,0,0,0,31.6467482405604805,35,'High','Microsoft 365','High','CLI-UPDATE',0,0,31.6467482405604805,60,'2026-04-03 12:20:32',80,'','High likelihood of reply from an established infrastructure.','ABANDONED','["+25: Tier-1 Infrastructure (Google/Microsoft)","+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(701,'loadtest-1@outlook.com','loadtest-1@outlook.com',0,0,1,1,0,0,0,31.6467482405604805,35,'High','Microsoft 365','High','CLI-UPDATE',0,0,31.6467482405604805,60,'2026-04-03 12:20:33',80,'','High likelihood of reply from an established infrastructure.','ABANDONED','["+25: Tier-1 Infrastructure (Google/Microsoft)","+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(702,'loadtest-21@outlook.com','loadtest-21@outlook.com',0,0,1,1,0,0,0,31.6467482405604805,35,'High','Microsoft 365','High','CLI-UPDATE',0,0,31.6467482405604805,60,'2026-04-03 12:20:33',80,'','High likelihood of reply from an established infrastructure.','ABANDONED','["+25: Tier-1 Infrastructure (Google/Microsoft)","+20: Established Legacy Identity (\u003e5 yrs)","+10: Consistent Historical Activity"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(703,'loadtest-13@protonmail.com','loadtest-13@protonmail.com',0,0,1,1,0,0,0,15.6269788899307275,35,'High','Proton Mail','High','CLI-UPDATE',0,0,15.6269788899307275,60,'2026-04-03 12:20:38',0,'550 5.1.1 <loadtest-13@protonmail.com>: Recipient address rejected: Address does not exist','This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(704,'loadtest-23@protonmail.com','loadtest-23@protonmail.com',0,0,1,1,0,0,0,15.6269788899307275,35,'High','Proton Mail','High','CLI-UPDATE',0,0,15.6269788899307275,60,'2026-04-03 12:20:39',0,'550 5.1.1 <loadtest-23@protonmail.com>: Recipient address rejected: Address does not exist','This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
INSERT INTO scans VALUES(705,'user@domain.com','user@domain.com',0,0,1,1,0,0,0,31.7763549748045655,35,'High','','High','CLI',0,0,31.7763549748045655,60,'2026-04-03 13:11:30',10,'550 5.7.1 Service unavailable, Client host [223.233.81.126] blocked using Spamhaus. To request removal from this list see https://www.spamhaus.org/query/ip/223.233.81.126 AS(1450) [DS2PEPF00003445.namprd04.prod.outlook.com 2026-04-03T13:11:28.370Z 08DE8F5643C8251D]','Low likelihood of reply due to infrastructure or identity age signals.','INVALID','["+20: Established Legacy Identity (\u003e5 yrs)","-15: Temporary Bounce/Greylisting detected","-30: SMTP Delivery Blocked (Security Filtering)"]','Reputation Warning: SMTP verify failed');
INSERT INTO scans VALUES(706,'legacy.user@gmail.com','legacy.user@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI',0,0,30.6584990413732718,100,'2026-04-03 13:11:42',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser 41be03b00d2f7-c76c648bd7esi10819887a12.32 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','INVALID','["Critical: Permanent Bounce (Engagement Zero)"]','Trusted: Tier-1 Infrastructure (Google)');
INSERT INTO scans VALUES(707,'support@company.com','support@company.com',0,0,1,0,0,0,1,31.757149264127463,0,'High','','High','CLI',0,0,31.757149264127463,60,'2026-04-03 13:11:48',0,'','This domain does not exist or has no mail servers configured.','INVALID','["Critical: Invalid Domain (DNS Failure)"]','');
INSERT INTO scans VALUES(708,'legacy.user@gmail.com','legacy.user@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI',0,0,30.6584990413732718,100,'2026-04-03 13:11:54',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser 41be03b00d2f7-c76c660bab7si12352834a12.344 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','INVALID','["Critical: Permanent Bounce (Engagement Zero)"]','Trusted: Tier-1 Infrastructure (Google)');
INSERT INTO scans VALUES(709,'sanjanamahajan2001@gmail.com','sanjanamahajan2001@gmail.com',0,1,1,1,1,0,0,30.6584990413732718,100,'Low','Google Workspace','High','CLI',0,0,30.6584990413732718,60,'2026-04-03 13:12:08',80,'250 2.1.5 OK','High likelihood of reply from an established infrastructure.','ACTIVE','["+25: Tier-1 Infrastructure (Google/Microsoft)","+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)"]','Trusted: Tier-1 Infrastructure (Google)');
INSERT INTO scans VALUES(710,'sanjanamaahi2201@gmail.com','sanjanamaahi2201@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI',0,0,30.6584990413732718,60,'2026-04-03 13:12:18',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser d9443c01a7336-2b2747390d4si109916045ad.22 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','INVALID','["Critical: Permanent Bounce (Engagement Zero)"]','Trusted: Tier-1 Infrastructure (Google)');
INSERT INTO scans VALUES(711,'sanjanamaahi2001@gmail.com','sanjanamaahi2001@gmail.com',0,1,1,1,1,0,0,30.6584990413732718,100,'Low','Google Workspace','High','CLI',0,0,30.6584990413732718,60,'2026-04-03 13:12:26',80,'250 2.1.5 OK','High likelihood of reply from an established infrastructure.','ACTIVE','["+25: Tier-1 Infrastructure (Google/Microsoft)","+20: Established Legacy Identity (\u003e5 yrs)","+10: Active Handshake (250 OK)"]','Trusted: Tier-1 Infrastructure (Google)');
INSERT INTO scans VALUES(712,'sanjana.mahajan@sofueled.com','sanjana.mahajan@sofueled.com',0,1,1,1,1,0,0,4.75434015314025693,100,'Low','','High','CLI',0,0,4.75434015314025693,60,'2026-04-03 13:12:41',45,'250 2.1.5 OK','Moderate likelihood. Common for enterprise or catch-all addresses.','ACTIVE','["+10: Active Handshake (250 OK)"]','Established Legacy Identity');
INSERT INTO scans VALUES(713,'abandoned-test@gmail.com',NULL,NULL,1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2026-03-03 13:12:47',0,NULL,NULL,'ACTIVE',NULL,NULL);
INSERT INTO scans VALUES(714,'abandoned-test@gmail.com','abandoned-test@gmail.com',0,0,1,1,0,0,0,30.6584990413732718,35,'High','Google Workspace','High','CLI',0,0,30.6584990413732718,60,'2026-04-03 13:12:56',0,replace('550 5.1.1 The email account that you tried to reach does not exist. Please try\n5.1.1 double-checking the recipient''s email address for typos or\n5.1.1 unnecessary spaces. For more information, go to\n5.1.1  https://support.google.com/mail/?p=NoSuchUser 41be03b00d2f7-c76c65a750esi11162729a12.192 - gsmtp','\n',char(10)),'This address is unreachable (Hard Bounce). No engagement is possible.','ABANDONED','["Critical: Permanent Bounce (Engagement Zero)"]','Identity Abandoned: Previously Valid address now persistent failure');
CREATE TABLE disposable_domains (
		domain TEXT PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	, provider_name TEXT);
INSERT INTO disposable_domains VALUES('mailinator.com','2026-03-31 05:47:26',NULL);
INSERT INTO disposable_domains VALUES('sub.mailinator.com','2026-03-31 10:51:43','Mailinator Hub');
INSERT INTO disposable_domains VALUES('new.sub.mailinator.com','2026-03-31 12:21:09','Mailinator Hub');
INSERT INTO disposable_domains VALUES('guerrillamail.com','2026-03-31 12:24:33','GuerrillaMail Hub');
INSERT INTO disposable_domains VALUES('burner-mail-2026.com','2026-03-31 12:24:58','');
INSERT INTO disposable_domains VALUES('atomicmail.io','2026-04-01 05:50:33','Confirmed via Active Probe: Found keyword ''temporary email''');
INSERT INTO disposable_domains VALUES('10minutemail.com','2026-04-02 04:37:07','');
CREATE TABLE metadata (
		key TEXT PRIMARY KEY,
		value TEXT
	);
INSERT INTO metadata VALUES('last_disposable_sync','2026-03-30T10:33:33Z');
INSERT INTO metadata VALUES('last_disposable_discovery','2026-04-02T09:56:34Z');
CREATE TABLE disposable_mx_signatures (
		signature TEXT PRIMARY KEY,
		provider_name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
INSERT INTO disposable_mx_signatures VALUES('mail.grr.la','GuerrillaMail Hub','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.dispostable.com','Dispostable Hub','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mail.mailinator.com','Mailinator Hub','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.emailondeck.com','EmailOnDeck','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mail.guerrillamail.com','GuerrillaMail Hub','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mail.mintemail.com','MintEmail','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.mailapi.org','Mail7','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.mail7.io','Mail7','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.yopmail.com','Yopmail','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.maildrop.cc','MailDrop Hub','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mail2.mailinator.com','Mailinator Hub','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.10minutemail.com','10MinuteMail','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('gate.poczta.onet.pl','Onet (Commonly used by PL-based disposables)','2026-04-02 11:24:15');
INSERT INTO disposable_mx_signatures VALUES('mx.disposable.com','Generic Disposable Hub','2026-04-02 11:24:15');
CREATE TABLE discovery_queue (
		domain TEXT PRIMARY KEY,
		mx_hosts TEXT,
		status TEXT DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
INSERT INTO discovery_queue VALUES('emailondeck.com','mail.protonmail.ch.,mailsec.protonmail.ch.','pending','2026-03-31 05:45:07');
INSERT INTO discovery_queue VALUES('github.com','github-com.mail.protection.outlook.com.','pending','2026-03-31 12:25:16');
INSERT INTO discovery_queue VALUES('gmail.com','gmail-smtp-in.l.google.com.,alt1.gmail-smtp-in.l.google.com.,alt2.gmail-smtp-in.l.google.com.,alt3.gmail-smtp-in.l.google.com.,alt4.gmail-smtp-in.l.google.com.','pending','2026-03-31 12:25:35');
INSERT INTO discovery_queue VALUES('microsoft.com','microsoft-com.mail.protection.outlook.com.','pending','2026-03-31 12:26:05');
INSERT INTO discovery_queue VALUES('sofueled.com','ASPMX.L.google.com.,ALT2.ASPMX.L.google.com.,ALT1.ASPMX.L.google.com.,ALT3.ASPMX.L.google.com.,ALT4.ASPMX.L.google.com.,rtsbu3a5fiert3y25qulvrmtvuoc4vahnc5fmlbms3hv6diid2ia.mx-verification.google.com.','pending','2026-04-01 05:10:53');
INSERT INTO discovery_queue VALUES('google.com','smtp.google.com.','pending','2026-04-01 05:52:46');
INSERT INTO discovery_queue VALUES('temp-mail.org','amir.mx.cloudflare.net.,linda.mx.cloudflare.net.,isaac.mx.cloudflare.net.','pending','2026-04-01 05:53:04');
INSERT INTO discovery_queue VALUES('amityonline.com','aspmx.l.google.com.,alt2.aspmx.l.google.com.,alt1.aspmx.l.google.com.,alt4.aspmx.l.google.com.,alt3.aspmx.l.google.com.','pending','2026-04-01 12:50:39');
INSERT INTO discovery_queue VALUES('tesla.com','tesla-com.mail.protection.outlook.com.','pending','2026-04-02 10:14:07');
INSERT INTO discovery_queue VALUES('stripe.com','aspmx.l.google.com.,alt1.aspmx.l.google.com.,alt2.aspmx.l.google.com.,aspmx2.googlemail.com.,aspmx3.googlemail.com.','pending','2026-04-02 11:01:50');
INSERT INTO discovery_queue VALUES('yopmail.com','smtp.yopmail.com.','pending','2026-04-02 11:02:36');
INSERT INTO discovery_queue VALUES('domain.com','["domain-com.mail.protection.outlook.com."]','PENDING','2026-04-02T12:04:35Z');
INSERT INTO discovery_queue VALUES('mailinator.com','["mail.mailinator.com.","mail2.mailinator.com."]','PENDING','2026-04-02T12:38:00Z');
INSERT INTO discovery_queue VALUES('example.com','["."]','PENDING','2026-04-03T05:15:20Z');
INSERT INTO discovery_queue VALUES('outlook.com','["outlook-com.olc.protection.outlook.com."]','PENDING','2026-04-03T05:30:52Z');
INSERT INTO discovery_queue VALUES('netflix.com','["aspmx.l.google.com.","alt2.aspmx.l.google.com.","alt1.aspmx.l.google.com.","aspmx2.googlemail.com.","aspmx3.googlemail.com."]','PENDING','2026-04-03T05:30:53Z');
INSERT INTO discovery_queue VALUES('yahoo.com','["mta6.am0.yahoodns.net.","mta7.am0.yahoodns.net.","mta5.am0.yahoodns.net."]','PENDING','2026-04-03T05:30:54Z');
INSERT INTO discovery_queue VALUES('zoho.com','["smtpin.zoho.com.","smtpin2.zoho.com.","smtpin3.zoho.com."]','PENDING','2026-04-03T05:30:54Z');
INSERT INTO discovery_queue VALUES('meta.com','["mxb-00082601.gslb.pphosted.com.","mxa-00082601.gslb.pphosted.com."]','PENDING','2026-04-03T05:30:56Z');
INSERT INTO discovery_queue VALUES('protonmail.com','["mail.protonmail.ch.","mailsec.protonmail.ch."]','PENDING','2026-04-03T05:31:04Z');
INSERT INTO discovery_queue VALUES('apple.com','["mx-in.g.apple.com.","mx-in-rn.apple.com.","mx-in-hfd.apple.com.","mx-in-sg.apple.com.","mx-in-vib.apple.com.","mx-in-ma.apple.com."]','PENDING','2026-04-03T05:31:43Z');
INSERT INTO discovery_queue VALUES('mail.com','["mx00.mail.com.","mx01.mail.com."]','PENDING','2026-04-03T05:34:47Z');
INSERT INTO discovery_queue VALUES('yandex.com','["mx.yandex.ru."]','PENDING','2026-04-03T05:34:47Z');
INSERT INTO discovery_queue VALUES('gmx.com','["mx01.gmx.net.","mx00.gmx.net."]','PENDING','2026-04-03T05:34:48Z');
INSERT INTO discovery_queue VALUES('twitter.com','["aspmx.l.google.com.","alt1.aspmx.l.google.com.","alt2.aspmx.l.google.com.","alt4.aspmx.l.google.com.","alt3.aspmx.l.google.com."]','PENDING','2026-04-03T05:34:48Z');
INSERT INTO discovery_queue VALUES('icloud.com','["mx01.mail.icloud.com.","mx02.mail.icloud.com."]','PENDING','2026-04-03T05:34:48Z');
INSERT INTO discovery_queue VALUES('fastmail.com','["in1-smtp.messagingengine.com.","in2-smtp.messagingengine.com."]','PENDING','2026-04-03T05:34:50Z');
INSERT INTO discovery_queue VALUES('facebook.com','["smtpin.vvv.facebook.com."]','PENDING','2026-04-03T05:34:53Z');
CREATE TABLE mx_signatures (
			signature TEXT PRIMARY KEY,
			provider_name TEXT,
			created_at DATETIME
		);
INSERT INTO mx_signatures VALUES('mx.emailondeck.com','EmailOnDeck','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.10minutemail.com','10MinuteMail','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mail.mintemail.com','MintEmail','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.mail7.io','Mail7','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.yopmail.com','Yopmail','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.disposable.com','Generic Disposable Hub','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mail.guerrillamail.com','GuerrillaMail Hub','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.mailapi.org','Mail7','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('gate.poczta.onet.pl','Onet (Commonly used by PL-based disposables)','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mail.grr.la','GuerrillaMail Hub','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.maildrop.cc','MailDrop Hub','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mx.dispostable.com','Dispostable Hub','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mail.mailinator.com','Mailinator Hub','2026-04-03T13:12:53Z');
INSERT INTO mx_signatures VALUES('mail2.mailinator.com','Mailinator Hub','2026-04-03T13:12:53Z');
CREATE TABLE domain_intelligence (
			domain TEXT PRIMARY KEY,
			age REAL,
			provider TEXT,
			trust TEXT,
			created_at DATETIME
		);
INSERT INTO domain_intelligence VALUES('gmail.com',30.6584990413732718,'Google Workspace','High','2026-04-02 12:27:05.769526323+00:00');
INSERT INTO domain_intelligence VALUES('stripe.com',30.5763264203471365,'Tier-1 Enterprise','High','2026-04-02 12:37:09.993071175+00:00');
INSERT INTO domain_intelligence VALUES('company.com',31.757149264127463,'','High','2026-04-02 12:37:39.19369523+00:00');
INSERT INTO domain_intelligence VALUES('mailinator.com',22.7669846330977954,'Custom/Unknown','High','2026-04-02 12:38:00.395461176+00:00');
INSERT INTO domain_intelligence VALUES('domain.com',31.7763549748045655,'','High','2026-04-02 12:52:10.48602743+00:00');
INSERT INTO domain_intelligence VALUES('example.com',30.6576775953303801,'','High','2026-04-03 05:15:20.646521549+00:00');
INSERT INTO domain_intelligence VALUES('outlook.com',31.6467482405604805,'Microsoft 365','High','2026-04-03 05:30:52.514571676+00:00');
INSERT INTO domain_intelligence VALUES('netflix.com',28.4110176760871908,'','High','2026-04-03 05:30:53.433308937+00:00');
INSERT INTO domain_intelligence VALUES('microsoft.com',34.9453784174841359,'Microsoft 365','High','2026-04-03 05:30:53.774027008+00:00');
INSERT INTO domain_intelligence VALUES('zoho.com',22.2258479015073788,'Zoho Workplace','High','2026-04-03 05:30:54.422098662+00:00');
INSERT INTO domain_intelligence VALUES('yahoo.com',31.2274560669390162,'Global Consumer','High','2026-04-03 05:30:54.527106617+00:00');
INSERT INTO domain_intelligence VALUES('github.com',18.4944255706368103,'Tier-1 Enterprise','High','2026-04-03 05:30:54.795905167+00:00');
INSERT INTO domain_intelligence VALUES('meta.com',35.2219766763897439,'','High','2026-04-03 05:30:56.466825851+00:00');
INSERT INTO domain_intelligence VALUES('google.com',28.56729635540254,'Google Workspace','High','2026-04-03 05:30:57.864178439+00:00');
INSERT INTO domain_intelligence VALUES('protonmail.com',15.6269788899307275,'Proton Mail','High','2026-04-03 05:31:04.273054404+00:00');
INSERT INTO domain_intelligence VALUES('apple.com',39.145265854102881,'Global Consumer','High','2026-04-03 05:31:43.975276922+00:00');
INSERT INTO domain_intelligence VALUES('mail.com',29.046641539926135,'','High','2026-04-03 05:34:47.603276465+00:00');
INSERT INTO domain_intelligence VALUES('gmx.com',31.9289475005568306,'','High','2026-04-03 05:34:48.377842642+00:00');
INSERT INTO domain_intelligence VALUES('yandex.com',27.5426461316319262,'','High','2026-04-03 05:34:48.407416745+00:00');
INSERT INTO domain_intelligence VALUES('twitter.com',26.2151950631222341,'','High','2026-04-03 05:34:48.510775169+00:00');
INSERT INTO domain_intelligence VALUES('icloud.com',27.2329429458297482,'Global Consumer','High','2026-04-03 05:34:48.740001181+00:00');
INSERT INTO domain_intelligence VALUES('fastmail.com',31.3370525941345796,'','High','2026-04-03 05:34:50.608807117+00:00');
INSERT INTO domain_intelligence VALUES('facebook.com',29.0329430923180673,'','High','2026-04-03 05:34:53.359471096+00:00');
INSERT INTO domain_intelligence VALUES('sofueled.com',4.75434015314025693,'','High','2026-04-03 13:12:41.069591043+00:00');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('scans',714);
CREATE INDEX idx_scans_email ON scans(email);
CREATE INDEX idx_scans_created ON scans(created_at);
CREATE INDEX idx_discovery_status ON discovery_queue(status);
COMMIT;
