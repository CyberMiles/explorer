import React, { Component } from "react"
import { Switch, Route } from "react-router-dom"
import { inject, observer } from "mobx-react"

import { Container, Segment, Label } from "semantic-ui-react"
import Account from "../Account"

@inject("commonStore")
@observer
export default class Home extends Component {
  componentDidMount() {
    this.props.commonStore.loadStatus()
  }
  componentWillReceiveProps(nextProps) {
    this.props.commonStore.loadStatus()
  }
  render() {
    const { isLoading, status } = this.props.commonStore
    return (
      <Container style={{ marginTop: "6em" }}>
        <Container style={{ marginBottom: "1em" }}>
          <Segment basic compact loading={isLoading}>
            <Label>
              LAST BLOCK:
              <Label.Detail>{status.latest_block_height}</Label.Detail>
            </Label>
          </Segment>
        </Container>
        <Switch>
          <Route path="/account/:address" component={Account} />
        </Switch>
      </Container>
    )
  }
}
