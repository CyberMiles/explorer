import { checkStatus, parseJSON, API_ROOT } from "./utils"

function get(height) {
  return fetch(`${API_ROOT}/block/${height}`)
    .then(checkStatus)
    .then(parseJSON)
}

export default { get }
