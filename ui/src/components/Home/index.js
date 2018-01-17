import React, { Component } from "react"
import { Switch, Route, Link } from "react-router-dom"
import { inject, observer } from "mobx-react"
import { Container, Segment, Label } from "semantic-ui-react"

import Account from "../Account"
import Block from "../Block"
import Transaction from "../Transaction"

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
              <Label.Detail as={Link} to={"/block/" + status.latest_block_height}>
                {status.latest_block_height}
              </Label.Detail>
            </Label>
          </Segment>
        </Container>
        <Switch>
          <Route path="/account/:address" component={Account} />
          <Route path="/block/:height" component={Block} />
          <Route path="/tx/:txhash" component={Transaction} />
        </Switch>
      </Container>
    )
  }
}
