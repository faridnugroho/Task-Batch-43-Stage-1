let projects = [];

function getData(event) {
  event.preventDefault();

  let name = document.getElementById("name").value;
  let start_date = document.getElementById("start_date").value;
  let end_date = document.getElementById("end_date").value;
  let description = document.getElementById("description").value;
  let image = document.getElementById("image").files;
  let nodeJS = document.getElementById("nodejs");
  let reactJS = document.getElementById("reactjs");
  let nextJS = document.getElementById("nextjs");
  let typeScript = document.getElementById("typescript");

  image = URL.createObjectURL(image[0]);

  //Deklarasi variable buat icon technology
  let nodejs = "";
  let reactjs = "";
  let nextjs = "";
  let typescript = "";

  //Pengkondisian buat masukin img icon ke variable
  nodeJS.checked == true ? (nodejs = "icons8-node-js") : "";
  reactJS.checked == true ? (reactjs = "icons8-react") : "";
  nextJS.checked == true ? (nextjs = "icons8-next-js") : "";
  typeScript.checked == true ? (typescript = "icons8-typescript") : "";

  let addProject = {
    name,
    start_date,
    end_date,
    description,
    image,
    nodejs,
    reactjs,
    nextjs,
    typescript,
  };

  projects.push(addProject);

  showData();
}

function showData() {
  document.getElementById("content").innerHTML = "";

  for (let i = 0; i < projects.length; i++) {
    document.getElementById("content").innerHTML += `
        <div class="project-card" id="project-card">
            <img src="${projects[i].image}" alt="project">
            <h4><a href="./project-detail.html" target="_blank">${
              projects[i].name
            }</a></h4>
            <h5>${getDistanceTime(
              projects[i].start_date,
              projects[i].end_date
            )}</h5>
            <p>${projects[i].description}</p>
            <div class="project-icon">
              <div class="${projects[i].nodejs}"></div>
              <div class="${projects[i].reactjs}"></div>
              <div class="${projects[i].nextjs}"></div>
              <div class="${projects[i].typescript}"></div>
            </div>
            <div class="project-button">
                <button>edit</button>
                <button>delete</button>
            </div>
        </div>
    `;
  }
}

function getDistanceTime(start_date, end_date) {
  let timeNow = new Date(start_date);
  let timeEnd = new Date(end_date);

  let distance = timeEnd - timeNow;

  let distanceDay = Math.floor(distance / (1000 * 60 * 60 * 24));
  let distanceHours = Math.floor(distance / (1000 * 60 * 60));
  let distanceMinutes = Math.floor(distance / (1000 * 60));
  let distanceSecond = Math.floor(distance / 1000);

  if (distanceDay > 0) {
    return `${distanceDay} Day`;
  } else if (distanceHours > 0) {
    return `${distanceHours} Hours`;
  } else if (distanceMinutes > 0) {
    return `${distanceMinutes} Minutes`;
  } else {
    return `${distanceSecond} Second`;
  }
}
