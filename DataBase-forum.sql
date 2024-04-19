-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Hôte : 127.0.0.1:3306
-- Généré le : ven. 19 avr. 2024 à 14:00
-- Version du serveur : 8.0.31
-- Version de PHP : 8.0.26

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de données : `forum`
--

-- --------------------------------------------------------

--
-- Structure de la table `messages`
--

DROP TABLE IF EXISTS `messages`;
CREATE TABLE IF NOT EXISTS `messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `contenu` text,
  `auteur` varchar(255) DEFAULT NULL,
  `date_creation` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `sujet_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `auteur` (`auteur`(250)),
  KEY `sujet_id` (`sujet_id`)
) ENGINE=MyISAM AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Déchargement des données de la table `messages`
--

INSERT INTO `messages` (`id`, `contenu`, `auteur`, `date_creation`, `sujet_id`) VALUES
(102, 'er', 'Satore', '2024-04-16 08:01:10', 55),
(101, 'cc', 'Satore', '2024-04-16 07:58:49', 55),
(100, 'aa', 'Natounor', '2024-04-16 07:53:47', 55),
(99, 'yo', 'Natounor', '2024-04-08 07:38:41', 45),
(98, 'dd', 'Satore', '2024-04-08 07:38:02', 45),
(96, 'nkk', 'Satore', '2024-03-25 11:09:26', 43),
(95, 'yo', 'Satore', '2024-03-25 10:15:00', 30),
(87, 'CC', 'Natounor', '2024-03-25 09:29:52', 43),
(88, 'CS', 'Natounor', '2024-03-25 09:29:54', 43),
(89, 'yo', 'Satore', '2024-03-25 09:30:20', 43),
(90, 'dz', 'Satore', '2024-03-25 09:30:22', 43),
(91, 'dz', 'Satore', '2024-03-25 09:30:23', 43),
(92, 'dz', 'Satore', '2024-03-25 09:30:24', 43),
(93, 'dz', 'Satore', '2024-03-25 09:30:25', 43),
(94, 'dz', 'Satore', '2024-03-25 09:30:27', 43);

-- --------------------------------------------------------

--
-- Structure de la table `sujets`
--

DROP TABLE IF EXISTS `sujets`;
CREATE TABLE IF NOT EXISTS `sujets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `titre` varchar(255) DEFAULT NULL,
  `contenu` text,
  `auteur` varchar(255) DEFAULT NULL,
  `date_creation` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `nomDuJeux` varchar(255) DEFAULT NULL,
  `like` int DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `auteur` (`auteur`(250))
) ENGINE=MyISAM AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Déchargement des données de la table `sujets`
--

INSERT INTO `sujets` (`id`, `titre`, `contenu`, `auteur`, `date_creation`, `nomDuJeux`, `like`) VALUES
(30, 'ouaii', 'cc', 'Satore', '2024-03-15 09:38:13', 'The witcher', 0),
(43, 'Test Tchat', 'cc', 'Natounor', '2024-03-22 11:20:45', 'Final Fantasy', 0),
(44, 'Test2', 'tom ', 'Satore', '2024-03-25 08:00:38', 'Final Fantasy', 0),
(45, 'Test', 'ddzd', 'Satore', '2024-03-25 10:22:46', 'Zelda', 0),
(55, 'Test', 'Test', 'Natounor', '2024-04-16 07:20:42', 'Yakuza', 0),
(56, 'Test2', 'Test2', 'Natounor', '2024-04-16 07:22:12', 'Yakuza', 0),
(52, 'Test', 'zaza', 'Natounor', '2024-04-08 08:05:51', 'The witcher', 0);

-- --------------------------------------------------------

--
-- Structure de la table `utilisateurs`
--

DROP TABLE IF EXISTS `utilisateurs`;
CREATE TABLE IF NOT EXISTS `utilisateurs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nom` varchar(50) DEFAULT NULL,
  `prenom` varchar(50) DEFAULT NULL,
  `pseudo` varchar(50) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `mot_de_passe` varchar(100) DEFAULT NULL,
  `image_profil` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pseudo` (`pseudo`),
  UNIQUE KEY `email` (`email`)
) ENGINE=MyISAM AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Déchargement des données de la table `utilisateurs`
--

INSERT INTO `utilisateurs` (`id`, `nom`, `prenom`, `pseudo`, `email`, `mot_de_passe`, `image_profil`) VALUES
(23, 'Novarese', 'Nathan', 'Natounor', 'nathan@ynov.com', 'nathan', '56e14a35-0579-484a-8e9c-b0501184005e.jpg'),
(22, 'Ballauri', 'Tom', 'Satore', 'tom@ynov.com', 'tom', '71eb011f-6e39-48ea-a39a-c88a831c8f30.jpg'),
(20, 'AZA', 'AZA', 'AZA', 'AZA@ynov.com', 'AZA', 'ea4f0295-dc77-42f4-b5b4-e5202fa66bda.jpg');
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
