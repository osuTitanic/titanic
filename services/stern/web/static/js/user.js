// The active tab is derived from the url #<tab>, defaulting to the general tab.
// All other tabs are fetched on user interaction from their HTML partials.

var activeTab = window.location.hash !== "" ? window.location.hash.replace("#", "") : "general";

function loadTab(id) {
    var content = document.getElementById(id);

    if (!content) return;
    if (content.getAttribute("data-loaded") === "true") return;

    var partial = content.getAttribute("data-partial");
    if (!partial) return;

    // Mark the tab as loaded before the request, to prevent duplicate requests
    content.setAttribute("data-loaded", "true");
    $(content).load(partial + "?mode=" + mode);
}

function expandProfileTab(id, forceExpand) {
    var tab = document.getElementById(id);

    if (!tab) {
        expandProfileTab("general", forceExpand);
        return;
    }

    activeTab = id;
    var isExpanded = tab.className.indexOf("expanded") !== -1;

    // If forceExpand is false and the tab is already expanded, collapse it
    if (!forceExpand && isExpanded) {
        $(tab).removeClass("expanded");
        slideUp(tab);
        return;
    }

    // Otherwise, expand it
    loadTab(id);

    if ($(tab).is(":hidden") || tab.style.height === "0px") {
        slideDown(tab);
        setTimeout(function () {
            tab.className += " expanded";
        }, 500);
    }

    if (forceExpand) {
        window.location.hash = "#" + activeTab;
    }
}

function loadMoreActivity(link, offset) {
    var row = $(link).closest("tr");

    $.get("/partials/users/" + userId + "/activity?mode=" + mode + "&offset=" + offset, function (html) {
        row.replaceWith(html);
        renderTimeagoElements();
    });

    return false;
}

function updatePlaystyleElement(element) {
    var nowUsing = $(element).toggleClass("playstyle-using").hasClass("playstyle-using");
    performApiRequest(nowUsing ? "POST" : "DELETE", "/users/" + userId + "/playstyle", { playstyle: element.id });
}

function addFriend() {
    if (!isLoggedIn()) return false;

    performApiRequest("POST", "/account/friends?id=" + userId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var targetAdded = data.status === "mutual" || superFriendly;
        document.getElementById("friend-status").className = "friend-current-true-target-" + targetAdded;
    });

    return false;
}

function removeFriend() {
    if (!isLoggedIn()) return false;

    performApiRequest("DELETE", "/account/friends?id=" + userId, null, function (xhr) {
        var data = JSON.parse(xhr.responseText);
        var targetAdded = data.status === "mutual" || superFriendly;
        document.getElementById("friend-status").className = "friend-current-false-target-" + targetAdded;
    });

    return false;
}

$(document).ready(function () {
    expandProfileTab(activeTab);
});
