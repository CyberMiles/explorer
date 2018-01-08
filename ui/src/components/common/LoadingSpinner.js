import React from "react"
import { Loader } from "semantic-ui-react"

export default class LoadingSpinner extends React.Component {
  render() {
    return <Loader active inline="centered" />
  }
}
