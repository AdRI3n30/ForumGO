// Récupérer les éléments du slider
const slider = document.querySelector('.slider');
const slide3 = document.getElementById('slide-3');
const slide4 = document.getElementById('slide-4');

function toggleMenu(menuId) {
    const menuDeroulant = document.getElementById(menuId);
    menuDeroulant.style.display = (menuDeroulant.style.display === 'flex') ? 'none' : 'flex';
}


// Fonction pour ajuster l'opacité du slide 3 en fonction du slide 4
function handleSlideVisibility() {
    const sliderRect = slider.getBoundingClientRect();
    const slide4Rect = slide4.getBoundingClientRect();

    // Vérifier si le slide 4 est visible dans le slider
    if (
        slide4Rect.left >= sliderRect.left &&
        slide4Rect.right <= sliderRect.right
    ) {
        slide3.style.opacity = '0'; // Masquer le slide 3 lorsque le slide 4 est visible
    } else {
        slide3.style.opacity = '1'; // Afficher le slide 3 lorsque le slide 4 n'est pas visible
    }
}

// Écouter l'événement scroll pour mettre à jour la visibilité du slide 3
window.addEventListener('scroll', handleSlideVisibility);

// Appeler la fonction handleSlideVisibility() une fois au chargement de la page pour la configuration initiale
handleSlideVisibility();

// Fonction pour activer/désactiver le menu déroulant
function toggleMenu(menuId) {
    const menuDeroulant = document.getElementById(menuId);
    const allMenus = document.querySelectorAll('.menu-deroulant');

    // Masquer tous les autres menus déroulants
    allMenus.forEach(menu => {
        if (menu.id !== menuId) {
            menu.style.display = 'none';
        }
    });

    // Activer ou désactiver le menu déroulant cible
    menuDeroulant.style.display = (menuDeroulant.style.display === 'flex') ? 'none' : 'flex';
}

// Masquer le menu déroulant lorsqu'on clique à l'extérieur
document.addEventListener('click', function(event) {
    const allMenus = document.querySelectorAll('.menu-deroulant');

    allMenus.forEach(menu => {
        const menuButton = document.querySelector(`.${menu.id.replace('Menu', '')}`);
        if (!menuButton.contains(event.target) && !menu.contains(event.target)) {
            menu.style.display = 'none';
        }
    });
});

