import { checkStatus, parseJSON, API_ROOT } from "./utils"

function get() {
  return fetch(`${API_ROOT}/status`)
    .then(checkStatus)
    .then(parseJSON)
}

export default { get }
