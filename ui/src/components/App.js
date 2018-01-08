import Header from "./Header"
import React from "react"
import { Switch, Route, withRouter } from "react-router-dom"
import { inject, observer } from "mobx-react"

import Home from "./Home"

@inject("commonStore")
@withRouter
@observer
export default class App extends React.Component {
  componentWillMount() {
    if (!this.props.commonStore.token) {
      this.props.commonStore.setAppLoaded()
    }
  }

  render() {
    if (this.props.commonStore.appLoaded) {
      return (
        <div>
          <Header />
          <Switch>
            <Route path="/" component={Home} />
          </Switch>
        </div>
      )
    }
    return <Header />
  }
}
