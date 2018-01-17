import { checkStatus, parseJSON, API_ROOT } from "./utils"

function get(txhash) {
  return fetch(`${API_ROOT}/tx/${txhash}`)
    .then(checkStatus)
    .then(parseJSON)
}

export default { get }
