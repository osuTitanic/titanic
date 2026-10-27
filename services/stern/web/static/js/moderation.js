var defaultError = {
    error: 500,
    details: "An internal server error occurred."
};

function handleApiErrorCallback(xhr, handlerFunction) {
    var error = defaultError;
    try {
        error = JSON.parse(xhr.responseText);
    } catch (e) {}
    if (handlerFunction) {
        handlerFunction(error);
    }
}

function getUser(userId, onSuccess, onFailure) {
    var url = "/moderation/users/" + userId + "/profile";

    return performApiRequest(
        "GET",
        url,
        null,
        function (xhr) {
            var user = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(user);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function updateUserProfile(userId, data, onSuccess, onFailure) {
    return getUser(
        userId,
        function (user) {
            var profileUpdate = {
                country: user.country,
                email: user.email,
                is_bot: user.is_bot,
                activated: user.activated,
                discord_id: user.discord_id,
                userpage: user.userpage,
                signature: user.signature,
                title: user.title,
                banner: user.banner,
                website: user.website,
                discord: user.discord,
                twitter: user.twitter,
                location: user.location,
                interests: user.interests
            };

            // Update fields provided in "data"
            for (var key in data) {
                if (data.hasOwnProperty(key) && profileUpdate.hasOwnProperty(key)) {
                    profileUpdate[key] = data[key];
                }
            }

            var url = "/moderation/users/" + userId + "/profile";

            return performApiRequest(
                "PATCH",
                url,
                profileUpdate,
                function (xhr) {
                    var user = JSON.parse(xhr.responseText);
                    if (onSuccess) {
                        onSuccess(user);
                    }
                },
                function (xhr) {
                    var error = defaultError;
                    try {
                        error = JSON.parse(xhr.responseText);
                    } catch (e) {}
                    if (onFailure) {
                        onFailure(error);
                    }
                }
            );
        },
        function (error) {
            if (onFailure) {
                return onFailure(error);
            }
        }
    );
}

function removeUserAvatar(userId, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to remove this user's avatar?")) {
        return null;
    }

    return performApiRequest(
        "DELETE",
        "/moderation/users/" + userId + "/avatar",
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess();
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function addBadge(userId, badgeData, onSuccess, onFailure) {
    performApiRequest(
        "POST",
        "/moderation/users/" + userId + "/badges",
        badgeData,
        function (xhr) {
            var badge = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(badge);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function updateBadge(userId, badgeId, badgeData, onSuccess, onFailure) {
    performApiRequest(
        "PATCH",
        "/moderation/users/" + userId + "/badges/" + badgeId,
        badgeData,
        function (xhr) {
            var badge = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(badge);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function removeBadge(userId, badgeId, onSuccess, onFailure) {
    performApiRequest(
        "DELETE",
        "/moderation/users/" + userId + "/badges/" + badgeId,
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess();
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function addStamp(userId, stampData, onSuccess, onFailure) {
    performApiRequest(
        "POST",
        "/moderation/users/" + userId + "/stamps",
        stampData,
        function (xhr) {
            var stamp = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(stamp);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function updateStamp(userId, stampId, stampData, onSuccess, onFailure) {
    performApiRequest(
        "PATCH",
        "/moderation/users/" + userId + "/stamps/" + stampId,
        stampData,
        function (xhr) {
            var stamp = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(stamp);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function removeStamp(userId, stampId, onSuccess, onFailure) {
    performApiRequest(
        "DELETE",
        "/moderation/users/" + userId + "/stamps/" + stampId,
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess();
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function changeUserName(userId, name, onSuccess, onFailure) {
    performApiRequest(
        "POST",
        "/moderation/users/" + userId + "/names",
        { name: name },
        function (xhr) {
            if (onSuccess) {
                onSuccess(JSON.parse(xhr.responseText));
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function updateUserNameHistory(userId, nameId, nameData, onSuccess, onFailure) {
    performApiRequest(
        "PATCH",
        "/moderation/users/" + userId + "/names/" + nameId,
        nameData,
        function (xhr) {
            if (onSuccess) {
                onSuccess(JSON.parse(xhr.responseText));
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function deleteUserNameHistory(userId, nameId, onSuccess, onFailure) {
    performApiRequest(
        "DELETE",
        "/moderation/users/" + userId + "/names/" + nameId,
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess(JSON.parse(xhr.responseText));
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function createUserInfringement(userId, infringementData, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to create this infringement? (This may take a while)")) {
        return;
    }

    performApiRequest(
        "POST",
        "/moderation/users/" + userId + "/infringements",
        infringementData,
        function (xhr) {
            var infringement = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(infringement);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function updateUserInfringement(userId, infringementId, infringementData, onSuccess, onFailure) {
    performApiRequest(
        "PATCH",
        "/moderation/users/" + userId + "/infringements/" + infringementId,
        infringementData,
        function (xhr) {
            var infringement = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(infringement);
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function deleteUserInfringement(userId, infringementId, restoreScores, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to delete this infringement? (This may take a while)")) {
        return;
    }

    var url = "/moderation/users/" + userId + "/infringements/" + infringementId;
    if (restoreScores) {
        url += "?restore_scores=true";
    }

    performApiRequest(
        "DELETE",
        url,
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess();
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function wipeUserScores(userId, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to wipe this user's scores? (This may take a while)")) {
        return null;
    }
    if (!confirm("Are you ABSOLUTELY sure you want to WIPE the user's scores?")) {
        return null;
    }

    return performApiRequest(
        "DELETE",
        "/moderation/users/" + userId + "/scores",
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess();
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function restoreUserScores(userId, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to restore this user's scores? (This may take a while)")) {
        return null;
    }

    return performApiRequest(
        "POST",
        "/moderation/users/" + userId + "/scores/restore",
        null,
        function (xhr) {
            if (onSuccess) {
                onSuccess();
            }
        },
        function (xhr) {
            return handleApiErrorCallback(xhr, onFailure);
        }
    );
}

function deleteUserAccount(userId, onSuccess, onFailure) {
    // TODO
}

function clearUserProfile(userId, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to clear this user's profile? This action cannot be undone.")) {
        return null;
    }
    if (!confirm("Are you ABSOLUTELY sure you want to CLEAR the user's profile?")) {
        return null;
    }

    var data = {
        userpage: null,
        signature: null,
        title: null,
        banner: null,
        website: null,
        discord: null,
        twitter: null,
        location: null,
        interests: null
    };
    return updateUserProfile(userId, data, onSuccess, onFailure);
}

function updateUserCountry(userId, countryCode, onSuccess, onFailure) {
    var data = { country: countryCode.toUpperCase() };
    updateUserProfile(userId, data, onSuccess, onFailure);
}

function unlinkDiscord(userId, onSuccess, onFailure) {
    var data = { discord_id: null };
    updateUserProfile(userId, data, onSuccess, onFailure);
}

function setDiscord(userId, discordId, onSuccess, onFailure) {
    var data = { discord_id: discordId };
    updateUserProfile(userId, data, onSuccess, onFailure);
}

function setBotStatus(userId, isBot, onSuccess, onFailure) {
    var data = { is_bot: isBot };
    updateUserProfile(userId, data, onSuccess, onFailure);
}

function changeEmail(userId, newEmail, onSuccess, onFailure) {
    var data = { email: newEmail };
    updateUserProfile(userId, data, onSuccess, onFailure);
}

/* User Interface */

function defaultOnSuccess() {
    reloadPageSoon(250);
}

function defaultOnError(err) {
    alert("An error occurred: " + err.details);
}

var moderationLoaderOverlay = null;

function showModerationLoader() {
    if (!moderationLoaderOverlay) {
        moderationLoaderOverlay = createLoaderOverlay();
        document.body.appendChild(moderationLoaderOverlay);
    }
}

function hideModerationLoader() {
    if (moderationLoaderOverlay && moderationLoaderOverlay.parentNode) {
        moderationLoaderOverlay.parentNode.removeChild(moderationLoaderOverlay);
    }
    moderationLoaderOverlay = null;
}

function runModerationQuickAction(action) {
    var request = action(
        function () {
            defaultOnSuccess();
        },
        function (err) {
            hideModerationLoader();
            defaultOnError(err);
        }
    );

    if (request) {
        showModerationLoader();
    }
}

function moderationRemoveUserAvatar(userId) {
    runModerationQuickAction(function (onSuccess, onFailure) {
        return removeUserAvatar(userId, onSuccess, onFailure);
    });
}

function moderationClearUserProfile(userId) {
    runModerationQuickAction(function (onSuccess, onFailure) {
        return clearUserProfile(userId, onSuccess, onFailure);
    });
}

function moderationWipeUserScores(userId) {
    runModerationQuickAction(function (onSuccess, onFailure) {
        return wipeUserScores(userId, onSuccess, onFailure);
    });
}

function moderationRestoreUserScores(userId) {
    runModerationQuickAction(function (onSuccess, onFailure) {
        return restoreUserScores(userId, onSuccess, onFailure);
    });
}

function moderationSaveProfile(userId) {
    var data = {};

    var email = document.getElementById("mod-email").value.trim();
    if (email.length) data.email = email;

    var country = document.getElementById("mod-country").value.trim();
    if (country.length) data.country = country.toUpperCase();

    var discordId = document.getElementById("mod-discord-id").value.trim();
    data.discord_id = discordId.length ? discordId : null;
    data.title = document.getElementById("mod-title").value || null;
    data.website = document.getElementById("mod-website").value || null;
    data.twitter = document.getElementById("mod-twitter").value || null;
    data.discord = document.getElementById("mod-discord").value || null;
    data.location = document.getElementById("mod-location").value || null;
    data.interests = document.getElementById("mod-interests").value || null;
    data.userpage = document.getElementById("mod-userpage").value || null;
    data.signature = document.getElementById("mod-signature").value || null;
    data.is_bot = document.getElementById("mod-is-bot").checked;
    data.activated = document.getElementById("mod-activated").checked;

    updateUserProfile(
        userId,
        data,
        function (user) {
            document.getElementById("moderation-edit-profile").close();
            reloadPageSoon(250);
        },
        function (err) {
            alert("Error: " + err.details);
        }
    );
}

function moderationUpdateBadge(userId, badgeId) {
    var data = {
        badge_url: document.getElementById("badge-url-" + badgeId).value || null,
        icon_url: document.getElementById("badge-icon-" + badgeId).value || null,
        description: document.getElementById("badge-desc-" + badgeId).value || null
    };

    updateBadge(
        userId,
        badgeId,
        data,
        function (badge) {
            alert("Badge updated successfully!");
        },
        function (err) {
            alert("Failed to update badge: " + err.details);
        }
    );
}

function moderationDeleteBadge(userId, badgeId) {
    removeBadge(
        userId,
        badgeId,
        function () {
            var row = document.querySelector('tr[data-badge-id="' + badgeId + '"]');
            if (row) row.remove();
            else reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to delete badge: " + err.details);
        }
    );
}

function moderationAddBadge(userId) {
    var data = {
        badge_url: document.getElementById("badge-new-url").value || null,
        icon_url: document.getElementById("badge-new-icon").value || null,
        description: document.getElementById("badge-new-desc").value || null
    };

    addBadge(
        userId,
        data,
        function (badge) {
            // TODO: Insert new row for created badge
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to add badge: " + err.details);
        }
    );
}

function moderationUpdateStamp(userId, stampId) {
    var data = {
        stamp_url: document.getElementById("stamp-url-" + stampId).value || null,
        icon_url: document.getElementById("stamp-icon-" + stampId).value || null,
        description: document.getElementById("stamp-desc-" + stampId).value || null
    };

    updateStamp(
        userId,
        stampId,
        data,
        function (stamp) {
            alert("Stamp updated successfully!");
        },
        function (err) {
            alert("Failed to update stamp: " + err.details);
        }
    );
}

function moderationDeleteStamp(userId, stampId) {
    removeStamp(
        userId,
        stampId,
        function () {
            var row = document.querySelector('tr[data-stamp-id="' + stampId + '"]');
            if (row) row.remove();
            else reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to delete stamp: " + err.details);
        }
    );
}

function moderationAddStamp(userId) {
    var data = {
        stamp_url: document.getElementById("stamp-new-url").value || null,
        icon_url: document.getElementById("stamp-new-icon").value || null,
        description: document.getElementById("stamp-new-desc").value || null
    };

    addStamp(
        userId,
        data,
        function (stamp) {
            // TODO: Insert new row for created stamp
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to add stamp: " + err.details);
        }
    );
}

function moderationUpdateInfringement(userId, infringementId) {
    var data = {
        duration: (parseInt(document.getElementById("infringement-duration-" + infringementId).value, 10) || 0) * 60,
        description: document.getElementById("infringement-desc-" + infringementId).value || null,
        is_permanent: !!document.getElementById("infringement-permanent-" + infringementId).checked
    };

    updateUserInfringement(
        userId,
        infringementId,
        data,
        function (updated) {
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to update infringement: " + (err && err.details ? err.details : "Unknown error"));
        }
    );
}

function moderationDeleteInfringement(userId, infringementId, isRestriction) {
    var restoreScores = isRestriction ? confirm("Do you want to restore the user scores & stats?") : false;

    deleteUserInfringement(
        userId,
        infringementId,
        restoreScores,
        function () {
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to delete infringement: " + (err && err.details ? err.details : "Unknown error"));
        }
    );
}

function moderationAddInfringement(userId) {
    var action = parseInt(document.getElementById("infringement-new-action").value, 10) || 0;
    var duration = parseInt(document.getElementById("infringement-new-duration").value, 10) || 0;
    var is_permanent = !!document.getElementById("infringement-new-permanent").checked;
    var description = document.getElementById("infringement-new-desc").value || null;

    var data = {
        duration: duration * 60,
        action: action,
        description: description,
        is_permanent: is_permanent
    };

    createUserInfringement(
        userId,
        data,
        function (created) {
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to add infringement: " + (err && err.details ? err.details : "Unknown error"));
        }
    );
}

function moderationChangeName(userId) {
    var name = document.getElementById("moderation-new-name").value.trim();
    if (!name) {
        alert("Enter a new name.");
        return;
    }
    if (!confirm("Change this user's name to " + name + "?")) {
        return;
    }

    changeUserName(
        userId,
        name,
        function () {
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to change name: " + err.details);
        }
    );
}

function moderationUpdateName(userId, nameId) {
    var data = {
        name: document.getElementById("name-value-" + nameId).value.trim(),
        reserved: document.getElementById("name-reserved-" + nameId).checked
    };
    if (!data.name) {
        alert("Previous name cannot be empty.");
        return;
    }

    updateUserNameHistory(
        userId,
        nameId,
        data,
        function () {
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to update name history: " + err.details);
        }
    );
}

function moderationDeleteName(userId, nameId) {
    deleteUserNameHistory(
        userId,
        nameId,
        function () {
            reloadPageSoon(250);
        },
        function (err) {
            alert("Failed to delete name history: " + err.details);
        }
    );
}
