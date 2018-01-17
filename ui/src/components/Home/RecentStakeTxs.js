import React, { Component } from "react"
import { Segment, Table, Header } from "semantic-ui-react"
import { inject, observer } from "mobx-react"
import { Link } from "react-router-dom"

import RedError from "../common/RedError"

@inject("txStore")
@observer
class RecentStakeTxs extends Component {
  componentDidMount() {
    this.props.txStore.loadRecentStakeTx()
  }
  render() {
    const { isLoading, error } = this.props.txStore
    if (error !== undefined) return <RedError message={error} />

    const recent = this.props.txStore.recentStakeTx || []
    return (
      <Segment basic loading={isLoading}>
        <Header>Stake Trans</Header>
        <Table compact fixed singleLine>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Hash</Table.HeaderCell>
              <Table.HeaderCell>Type</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {recent.map((v, index) => {
              return (
                <Table.Row key={index} verticalAlign="top">
                  <Table.Cell>
                    <Link to={"/tx/" + v.txhash}>{v.txhash}</Link>
                  </Table.Cell>
                  <Table.Cell>{v.type}</Table.Cell>
                </Table.Row>
              )
            })}
          </Table.Body>
        </Table>
      </Segment>
    )
  }
}

export default RecentStakeTxs
