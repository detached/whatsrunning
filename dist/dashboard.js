function createStageNode(project, deployment) {
    var stageElement = document.createElement('div');
    stageElement.id = createId(project, deployment.stage);
    stageElement.appendChild(createHeader(deployment.stage, 'stage-header'));
    stageElement.appendChild(createContent(createId(stageElement.id, "version"), deployment.version));
    return stageElement;
}

function createProjectNode(project) {
    var projectElement = document.createElement('div');
    projectElement.id = project.name;
    projectElement.className = 'project';
    projectElement.appendChild(createHeader(project.name, 'project-header'));

    project.deployments.forEach(function (deployment) {
        projectElement.appendChild(createStageNode(project.name, deployment));
    });

    return projectElement;
}

function renderProjects() {
    var projectsScriptTag = document.getElementById('projects');
    var rootElement = document.getElementById('root');
    var projects = JSON.parse(projectsScriptTag.innerHTML);

    if (projects) {
        projects.forEach(function (project) {
            rootElement.appendChild(createProjectNode(project));
        })
    }
}

function createNotification(text) {
    var notification = document.createElement('div');
    notification.className = 'notification';
    notification.appendChild(document.createTextNode(text));
    return notification;
}

function listenForChanges() {

    var connection = new WebSocket('ws://' + document.location.host + '/ws');
    connection.onmessage = function(e) {
        var message = JSON.parse(e.data);

        console.log(message);
        if (message.action === 'update') {
            addOrUpdate(message.project);
        }
    };

    connection.onclose = function (e) {
        var body = document.createElement('body');
        body.appendChild(createNotification('Connection to server lost'));
        console.log('Connection closed:', JSON.stringify(e))
    };
}

function updateProjectNode(projectNode, project) {
    project.deployments.forEach(function (deployment) {

        var stageId = createId(project.name, deployment.stage);
        var stageNode = document.getElementById(stageId);

        if (stageNode) {

            var versionId = createId(stageId, 'version');
            stageNode.replaceChild(createContent(versionId, deployment.version), document.getElementById(versionId));
        } else {

            projectNode.appendChild(createStageNode(project.name, deployment));
        }
    });
}

function addOrUpdate(project) {

    var projectNode = document.getElementById(project.name);

    if (projectNode) {
        updateProjectNode(projectNode, project);
    } else {
        projectNode = createProjectNode(project);
        var rootElement = document.getElementById('root');
        rootElement.appendChild(projectNode);
    }
}

function createHeader(name, cssClass) {
    var header = document.createElement('span');
    header.appendChild(document.createTextNode(name));
    header.className = cssClass;
    return header;
}

function createContent(id, text) {
    var content = document.createElement('p');
    content.id = id;
    content.className = 'content';
    content.appendChild(document.createTextNode(text));
    return content;
}

function createId(project, stage) {
    return project + "-" + stage ;
}