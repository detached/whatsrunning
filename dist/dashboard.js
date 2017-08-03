function createProjectNode(project) {
    var projectElement = document.createElement('div');
    projectElement.id = project.name;
    projectElement.className = 'project';
    projectElement.appendChild(createHeader(project.name, 'project-header'));

    project.deployments.forEach(function (deployment) {

        var stageElement = document.createElement('div');
        stageElement.id = deployment.stage;
        stageElement.appendChild(createHeader(deployment.stage, 'stage-header'));
        stageElement.appendChild(createContent("version", deployment.version));

        projectElement.appendChild(stageElement);
    });

    return projectElement;
}

function renderProjects() {
    var projectsScriptTag = document.getElementById('projects');
    var rootElement = document.getElementById('root');
    var projects = JSON.parse(projectsScriptTag.innerHTML);

    projects.forEach(function (project) {
        rootElement.appendChild(createProjectNode(project));
    })
}

function listenForChanges() {

    var connection = new WebSocket("ws://" + document.location.host + "/ws");
    connection.onmessage = function(e) {
        var message = e.data;

        if (message.action === 'update') {
            addOrUpdate(message.project);
        }
    };
    connection.onclose = function (e) {
        console.log("Connection closed:", JSON.stringify(e))
    };
}

function updateProjectNode(projectNode, project) {
    project.deployments.forEach(function (deployment) {

        var stageNode = projectNode.getElementById(deployment.stage);

        if (stageNode) {
            stageNode.set
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