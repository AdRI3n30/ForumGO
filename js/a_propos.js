document.addEventListener('DOMContentLoaded', function () {
    var menuButtonJeux = document.getElementById('menuButtonJeux');
    var menuJeux = document.getElementById('menuJeux');

    menuButtonJeux.addEventListener('click', function () {
        toggleMenu(menuJeux);
    });

    var menuButtonTopics = document.getElementById('menuButtonTopics');
    var menuTopics = document.getElementById('menuTopics');

    menuButtonTopics.addEventListener('click', function () {
        toggleMenu(menuTopics);
    });

    function toggleMenu(menu) {
        menu.classList.toggle('show'); // Ajoute ou supprime la classe 'show'
    }

    // Fermer le menu si l'utilisateur clique en dehors du menu
    document.addEventListener('click', function (event) {
        if (!menuJeux.contains(event.target) && event.target !== menuButtonJeux && !menuTopics.contains(event.target) && event.target !== menuButtonTopics) {
            menuJeux.classList.remove('show'); // Assurez-vous que la classe est supprimée si on clique en dehors du menu Jeux
            menuTopics.classList.remove('show'); // Assurez-vous que la classe est supprimée si on clique en dehors du menu Topics
        }
    });
});
