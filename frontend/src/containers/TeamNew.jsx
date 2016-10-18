import { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class TeamNew extends Component {
  static propTypes = {
    createTeam: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleSuccess = (x) => {
    browserHistory.push(`/teams/${x.id}`);
    this.props.messagePush('created a team');
  };

  render() {
    return (
      <New
        resourceName="team"
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          logo: { name: 'Square logo image' },
        }}
        createResource={this.props.createTeam}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
