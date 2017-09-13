import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';

import Settings from './Settings';

const teamIdSelector = (state, props) => props.params.id;
const teamSelector = createSelector(
  selectors.teams,
  teamIdSelector,
  utils.get,
);

const membersSelector = createSelector(
  selectors.userTeams,
  teamIdSelector,
  (userTeams, teamId) =>
    userTeams.toList().filter(x =>
      x.get('teamId') === teamId &&
      x.get('leftAt') === null,
    ),
);

const myMembershipSelector = createSelector(
  membersSelector,
  selectors.myId,
  (members, myId) => members.find(x => x.get('userId') === myId),
);

@connect(createStructuredSelector({
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  teamId: teamIdSelector,
  team: teamSelector,
  members: membersSelector,
  myMembership: myMembershipSelector,
}), actions)
export default class TeamSettings extends Component {
  static propTypes = {
    myId: PropTypes.string,
    amAdmin: PropTypes.bool,
    teamId: PropTypes.string.isRequired,
    team: ImmutablePropTypes.map,
    members: ImmutablePropTypes.list,
    myMembership: ImmutablePropTypes.map,

    messagePush: PropTypes.func.isRequired,
    loadTeam: PropTypes.func.isRequired,
    loadUserTeams: PropTypes.func.isRequired,
    updateTeam: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ myId, amAdmin, teamId, team, members, myMembership }, prevProps) {
    if (!myId) {
      this.props.messagePush(
        'you need to be logged in to change team settings', true,
      );
      browserHistory.push('/login');
      return;
    }

    if (
      !prevProps ||
      teamId !== prevProps.teamId ||
      (!team && prevProps.team)
    ) {
      this.props.loadTeam(teamId);
      this.props.loadUserTeams(
        { teamId, leftAt: '%00' },
        members.reduce((r, x) => r.concat(x.get('id')), []),
      );
    }

    if (!team || (!myMembership && !amAdmin)) return;
    if (!amAdmin && !myMembership.get('isLeader')) {
      this.props.messagePush(
        'you need to be a leader or an admin to change team settings', true,
      );
      browserHistory.push('/');
    }
  }

  render() {
    const { team, updateTeam } = this.props;
    return (
      <Settings
        resourceName="team"
        resource={team}
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Tag / Abbreviation' },
          logo: { name: 'Square logo image' },
        }}
        updateResource={updateTeam}
      />
    );
  }
}
