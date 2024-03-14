var activeMenuId = null;

var activeMenuId = null;

function toggleMenu(menuId, triangleClass, textId, subMenuId, subTextId) {
    var triangle = document.querySelector("." + triangleClass);
    var menu = document.getElementById(menuId);
    var text = document.getElementById(textId);
    var subMenu = document.getElementById(subMenuId);
    var subText = document.getElementById(subTextId);

    // Si le menu est actuellement ouvert, ferme-le
    if (activeMenuId === menuId) {
        hideMenu(menuId);
        hideSubMenu(subMenuId, subTextId);
        return;
    }

    // Masquer tous les menus sauf celui actuel
    hideAllMenusExcept(menuId);

    if (menu.style.display === "none" || activeMenuId !== menuId) {
        triangle.style.transform = "rotate(90deg)";
        text.style.display = "block";
        menu.style.display = "block";
        activeMenuId = menuId;
    } else {
        hideMenu(menuId);
        hideSubMenu(subMenuId, subTextId);
    }
}

function hideAllMenusExcept(exceptMenuId) {
    // Liste des menus disponibles
    var menuIds = ["topicsMenu", "JeuxMenu"];

    // Masquer tous les menus sauf celui spécifié
    menuIds.forEach(function(menuId) {
        if (menuId !== exceptMenuId) {
            hideMenu(menuId);
        }
    });
}

function hideMenu(menuId) {
    var triangleClass = (menuId === "topicsMenu") ? "topics-triangle" : "jeux-triangle";
    var textId = (menuId === "topicsMenu") ? "topicsText" : "JeuxText";

    var triangle = document.querySelector("." + triangleClass);
    var menu = document.getElementById(menuId);
    var text = document.getElementById(textId);

    triangle.style.transform = "rotate(0deg)";
    text.style.display = "none";
    menu.style.display = "none";
    activeMenuId = null;
}

function toggleSubMenu(subMenuId, subTextId) {
    hideAllSubMenusExcept(subMenuId);

    var subMenu = document.getElementById(subMenuId);
    var subText = document.getElementById(subTextId);

    if (subMenu.style.display === "none") {
        subMenu.style.display = "block";
        subText.style.display = "block";
    } else {
        hideSubMenu(subMenuId, subTextId);
    }
}

function hideAllSubMenusExcept(exceptSubMenuId) {
    // Liste des sous-menus disponibles
    var subMenuIds = ["JeuxSubMenu"];

    // Masquer tous les sous-menus sauf celui spécifié
    subMenuIds.forEach(function(subMenuId) {
        if (subMenuId !== exceptSubMenuId) {
            hideSubMenu(subMenuId, subTextId);
        }
    });
}

function hideSubMenu(subMenuId, subTextId) {
    var subMenu = document.getElementById(subMenuId);
    var subText = document.getElementById(subTextId);

    subMenu.style.display = "none";
    subText.style.display = "none";
}

// Appel pour masquer tous les menus au chargement de la page
hideAllMenusExcept(null);
