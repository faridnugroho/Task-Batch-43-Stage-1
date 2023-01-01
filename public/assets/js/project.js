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

  // console.log(image);

  image = URL.createObjectURL(image[0]);

  //Deklarasi variable icon technology
  let nodejs = "";
  let reactjs = "";
  let nextjs = "";
  let typescript = "";

  //Pengkondisian icon ke variable
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
      <div class="col">
        <div class="card h-100">
          <img src="${
            projects[i].image
          }" class="card-img-top" alt="ImageProject">
          <div class="card-body">
            <h5 class="card-title">
              <a href="./project-detail.html" class="text-decoration-none text-dark" target="_blank">${
                projects[i].name
              }</a>
            </h5>
            <p class="text-light-emphasis">durasi : ${getDistanceTime(
              projects[i].start_date,
              projects[i].end_date
            )}</p>
            <p class="card-text me-4">${projects[i].description}</p>
            <div class="project-icon d-flex gap-1 my-4">
              <div class="${projects[i].nodejs}"></div>
              <div class="${projects[i].reactjs}"></div>
              <div class="${projects[i].nextjs}"></div>
              <div class="${projects[i].typescript}"></div>
            </div>
            <div class="d-flex gap-1">
              <a href="#" class="btn btn-sm btn-dark w-100">edit</a>
              <a href="#" class="btn btn-sm btn-dark w-100">delete</a>
            </div>
          </div>
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
