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

    performApiRequest(
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
    getUser(
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

            performApiRequest(
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
            return onFailure(error);
        }
    );
}

function removeUserAvatar(userId, onSuccess, onFailure) {
    if (!confirm("Are you sure you want to remove this user's avatar?")) {
        return;
    }

    performApiRequest(
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

function getUserInfringements(userId, onSuccess, onFailure) {
    performApiRequest(
        "GET",
        "/moderation/users/" + userId + "/infringements",
        null,
        function (xhr) {
            var infringements = JSON.parse(xhr.responseText);
            if (onSuccess) {
                onSuccess(infringements);
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
        return;
    }
    if (!confirm("Are you ABSOLUTELY sure you want to WIPE the user's scores?")) {
        return;
    }

    performApiRequest(
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
        return;
    }

    performApiRequest(
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
        return;
    }
    if (!confirm("Are you ABSOLUTELY sure you want to CLEAR the user's profile?")) {
        return;
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
    updateUserProfile(userId, data, onSuccess, onFailure);
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

function showSpinner() {
    var el = document.getElementById("moderation-loader");
    if (el) el.style.display = "block";
}

function hideSpinner() {
    var el = document.getElementById("moderation-loader");
    if (el) el.style.display = "none";
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

function moderationOpenInfringements(userId) {
    showSpinner();

    // Show dialog immediately & load data
    document.getElementById("moderation-infringements").showModal();
    var body = document.getElementById("moderation-infringements-body");

    // Remove existing rows except the new-row
    var rows = Array.from(body.querySelectorAll("tr"));
    $.each(rows, function (i, r) {
        if (r.id !== "infringement-new-row") r.remove();
    });

    getUserInfringements(
        userId,
        function (infringements) {
            infringements.sort(function (a, b) {
                return new Date(b.time) - new Date(a.time);
            });

            $.each(infringements, function (i, inf) {
                body.insertBefore(
                    createInfringementElement(inf, userId),
                    document.getElementById("infringement-new-row")
                );
            });
            hideSpinner();
        },
        function (err) {
            hideSpinner();
            alert("Failed to load infringements: " + (err && err.details ? err.details : "Unknown error"));
        }
    );
}

function moderationUpdateInfringement(userId, infringementId) {
    var row = document.querySelector('tr[data-infringement-id="' + infringementId + '"]');
    if (!row) return alert("Row not found");

    var data = {
        duration: (parseInt(row.elements.duration.value, 10) || 0) * 60,
        description: row.elements.description.value || null,
        is_permanent: !!row.elements.isPermanent.checked
    };
    showSpinner();

    updateUserInfringement(
        userId,
        infringementId,
        data,
        function (updated) {
            hideSpinner();
            moderationOpenInfringements(userId);
        },
        function (err) {
            hideSpinner();
            alert("Failed to update infringement: " + (err && err.details ? err.details : "Unknown error"));
        }
    );
}

function moderationDeleteInfringement(userId, infringementId, isRestriction) {
    var restoreScores = isRestriction ? confirm("Do you want to restore the user scores & stats?") : false;
    showSpinner();

    deleteUserInfringement(
        userId,
        infringementId,
        restoreScores,
        function () {
            hideSpinner();
            var row = document.querySelector('tr[data-infringement-id="' + infringementId + '"]');
            if (row) row.remove();
            else reloadPageSoon(250);
        },
        function (err) {
            hideSpinner();
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
    showSpinner();

    createUserInfringement(
        userId,
        data,
        function (created) {
            hideSpinner();
            moderationOpenInfringements(userId);
        },
        function (err) {
            hideSpinner();
            alert("Failed to add infringement: " + (err && err.details ? err.details : "Unknown error"));
        }
    );
}

function createInfringementElement(inf, userId) {
    var infringementRow = document.createElement("tr");
    infringementRow.setAttribute("data-infringement-id", inf.id);

    var timeTd = document.createElement("td");
    var timestamp = new Date(inf.time);
    timeTd.textContent = isNaN(timestamp.getTime()) ? "" : timestamp.toLocaleString();

    var actionTd = document.createElement("td");
    var moderationSelect = document.createElement("select");
    moderationSelect.style.width = "100%";

    var optionRestrict = document.createElement("option");
    optionRestrict.value = "0";
    optionRestrict.text = "Restrict";

    var optionSilence = document.createElement("option");
    optionSilence.value = "1";
    optionSilence.text = "Silence";

    moderationSelect.add(optionRestrict);
    moderationSelect.add(optionSilence);
    moderationSelect.value = String(inf.action);

    var durationTd = document.createElement("td");
    var durationInput = document.createElement("input");
    var durationEnd = new Date(inf.length);
    var length = Math.floor((durationEnd.getTime() - timestamp.getTime()) / 60000);
    durationInput.type = "number";
    durationInput.min = 0;
    durationInput.style.width = "100%";
    durationInput.value = Math.max(0, length);

    var isPermanentTd = document.createElement("td");
    isPermanentTd.style.textAlign = "center";
    var isPermanentCheckbox = document.createElement("input");
    isPermanentCheckbox.type = "checkbox";
    isPermanentCheckbox.checked = !!inf.is_permanent;

    var descriptionTd = document.createElement("td");
    var descriptionInput = document.createElement("input");
    descriptionInput.type = "text";
    descriptionInput.style.width = "100%";
    descriptionInput.value = inf.description || "";

    var actionsTd = document.createElement("td");
    actionsTd.style.textAlign = "right";

    var updBtn = document.createElement("button");
    updBtn.type = "button";
    updBtn.textContent = "Update";
    updBtn.onclick = function () {
        moderationUpdateInfringement(userId, inf.id);
    };

    var delBtn = document.createElement("button");
    delBtn.type = "button";
    delBtn.textContent = "Delete";
    delBtn.style.color = "#c00";
    delBtn.onclick = function () {
        if (confirm("Delete this infringement?")) moderationDeleteInfringement(userId, inf.id, inf.action === 0);
    };

    actionTd.appendChild(moderationSelect);
    durationTd.appendChild(durationInput);
    isPermanentTd.appendChild(isPermanentCheckbox);
    descriptionTd.appendChild(descriptionInput);
    actionsTd.appendChild(updBtn);
    actionsTd.appendChild(delBtn);

    infringementRow.appendChild(timeTd);
    infringementRow.appendChild(actionTd);
    infringementRow.appendChild(durationTd);
    infringementRow.appendChild(isPermanentTd);
    infringementRow.appendChild(descriptionTd);
    infringementRow.appendChild(actionsTd);

    // Store references for later retrieval
    infringementRow.elements = {
        action: moderationSelect,
        duration: durationInput,
        isPermanent: isPermanentCheckbox,
        description: descriptionInput
    };

    return infringementRow;
}
