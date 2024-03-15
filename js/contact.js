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

let map;

async function initMap() {
  const { Map } = await google.maps.importLibrary("maps");

  map = new Map(document.getElementById("map"), {
    center: { lat: 43.615171970344804, lng: 7.05530506998268},
    zoom: 18,
  });
}

initMap();