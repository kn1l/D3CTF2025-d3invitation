function generateInvitation(user_id, avatarFile) {
    if (avatarFile) {
        object_name = avatarFile.name;
        genSTSCreds(object_name)
            .then(credsData => {
                return putAvatar(
                    credsData.access_key_id,
                    credsData.secret_access_key,
                    credsData.session_token,
                    object_name,
                    avatarFile
                ).then(() => {
                    navigateToInvitation(
                        user_id,
                        credsData.access_key_id,
                        credsData.secret_access_key,
                        credsData.session_token,
                        object_name
                    )
                })
            })
            .catch(error => {
                console.error('Error generating STS credentials or uploading avatar:', error);
            });
    } else {
        navigateToInvitation(user_id);
    }
}


function navigateToInvitation(user_id, access_key_id, secret_access_key, session_token, object_name) {
    let url = `invitation?user_id=${encodeURIComponent(user_id)}`;

    if (access_key_id) {
        url += `&access_key_id=${encodeURIComponent(access_key_id)}`;
    }

    if (secret_access_key) {
        url += `&secret_access_key=${encodeURIComponent(secret_access_key)}`;
    }

    if (session_token) {
        url += `&session_token=${encodeURIComponent(session_token)}`;
    }

    if (object_name) {
        url += `&object_name=${encodeURIComponent(object_name)}`;
    }

    window.location.href = url;
}


function genSTSCreds(object_name) {
    return new Promise((resolve, reject) => {
        const genSTSJson = {
            "object_name": object_name
        }

        fetch('/api/genSTSCreds', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(genSTSJson)
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                resolve(data);
            })
            .catch(error => {
                reject(error);
            });
    });
}

function getAvatarUrl(access_key_id, secret_access_key, session_token, object_name) {
    return `/api/getObject?access_key_id=${encodeURIComponent(access_key_id)}&secret_access_key=${encodeURIComponent(secret_access_key)}&session_token=${encodeURIComponent(session_token)}&object_name=${encodeURIComponent(object_name)}`
}

function putAvatar(access_key_id, secret_access_key, session_token, object_name, avatar) {
    return new Promise((resolve, reject) => {
        const formData = new FormData();
        formData.append('access_key_id', access_key_id);
        formData.append('secret_access_key', secret_access_key);
        formData.append('session_token', session_token);
        formData.append('object_name', object_name);
        formData.append('avatar', avatar);

        fetch('/api/putObject', {
            method: 'POST',
            body: formData
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                resolve(data);
            })
            .catch(error => {
                reject(error);
            });
    });
}