import moment from 'moment';
import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import { Link } from 'react-router';

import * as actions from '../actions';
import * as selectors from '../selectors';
import * as utils from '../utils';
import Component from '../utils/Component';

import NotFound from './NotFound';
import TeamUserTeam from './TeamUserTeam';
import TeamUserTeamRequest from './TeamUserTeamRequest';

import Button from '../components/Button';

const teamIdSelector = (state, props) => props.params.id;
const teamSelector = createSelector(
  selectors.teams,
  teamIdSelector,
  utils.get,
);

const membersSelector = createSelector(
  selectors.userTeams,
  teamIdSelector,
  selectors.users,
  (userTeams, teamId, users) =>
    userTeams.toList().filter(x =>
      x.get('teamId') === teamId &&
      x.get('leftAt') === null
    ).sort((a, b) => {
      const leadership = +b.get('isLeader') - +a.get('isLeader');
      if (leadership !== 0) return leadership;
      const userA = users.get(a.get('userId'));
      const userB = users.get(b.get('userId'));
      if (!userA || !userB) return a.get('id') - b.get('id');
      return userA.get('nickname').localeCompare(userB.get('nickname'));
    }),
);

const myMembershipSelector = createSelector(
  membersSelector,
  selectors.myId,
  (members, myId) => members.find(x => x.get('userId') === myId),
);

const requestsSelector = createSelector(
  selectors.userTeamRequests,
  teamIdSelector,
  (userTeamRequests, teamId) =>
    userTeamRequests.toList().filter(x =>
      x.get('teamId') === teamId &&
      x.get('decision') === null
    ).sort((a, b) => a.get('id') - b.get('id')),
);

const myRequestSelector = createSelector(
  requestsSelector,
  selectors.myId,
  (requests, myId) => requests.find(x => x.get('userId') === myId),
);

const pastMembersSelector = createSelector(
  selectors.userTeams,
  teamIdSelector,
  (userTeams, teamId) =>
    userTeams.toList().filter(x =>
      x.get('teamId') === teamId &&
      x.get('leftAt') !== null
    ).sort((a, b) => a.get('id') - b.get('id')),
);

@connect(createStructuredSelector({
  myId: selectors.myId,
  amAdmin: selectors.amAdmin,
  users: selectors.users,

  teamId: teamIdSelector,
  team: teamSelector,
  members: membersSelector,
  myMembership: myMembershipSelector,
  requests: requestsSelector,
  myRequest: myRequestSelector,
  pastMembers: pastMembersSelector,
}), actions)
export default class Team extends Component {
  static propTypes = {
    myId: PropTypes.string,
    amAdmin: PropTypes.bool,
    users: ImmutablePropTypes.map,

    teamId: PropTypes.string,
    team: ImmutablePropTypes.map,
    members: ImmutablePropTypes.list,
    myMembership: ImmutablePropTypes.map,
    requests: ImmutablePropTypes.list,
    myRequest: ImmutablePropTypes.map,
    pastMembers: ImmutablePropTypes.list,

    loadTeam: PropTypes.func.isRequired,
    patchTeam: PropTypes.func.isRequired,
    createUserTeamRequest: PropTypes.func.isRequired,
    loadUserTeamRequests: PropTypes.func.isRequired,
    patchUserTeamRequest: PropTypes.func.isRequired,
    loadUserTeams: PropTypes.func.isRequired,
    patchUserTeam: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    if (
      nextProps.myId !== this.props.myId ||
      nextProps.amAdmin !== this.props.amAdmin ||
      nextProps.teamId !== this.props.teamId ||
      !utils.isSame(
        nextProps.myMembership, this.props.myMembership, ['isLeader'],
      )
    ) {
      this.load(nextProps);
    }
  }

  load({
    myId, amAdmin, teamId, members, myMembership, requests, pastMembers,
  }) {
    this.props.loadTeam(teamId);
    this.props.loadUserTeams(
      { teamId },
      members.concat(pastMembers).reduce((r, x) => r.concat(x.get('id')), []),
      ['userId'],
    );
    const requestIds = requests.reduce((r, x) => r.concat(x.get('id')), []);
    if ((myMembership && myMembership.get('isLeader')) || amAdmin) {
      this.props.loadUserTeamRequests(
        { teamId, decision: '%00' },
        requestIds,
        ['userId'],
      );
    } else if (myId && !myMembership) {
      this.props.loadUserTeamRequests(
        { userId: myId, teamId, decision: '%00' },
        requestIds,
      );
    }
  }

  handleLeave = () => {
    const { team } = this.props;
    if (!window.confirm(`Leave ${team && team.get('name')}?`)) return;
    this.props.patchUserTeam(this.props.myMembership.get('id'), 'leave', () =>
      this.props.messagePush('left a team'),
    );
  };

  handleAcceptInvite = () => {
    this.props.patchUserTeamRequest(
      this.props.myRequest.get('id'), 'yes', (x) => {
        x.decision && this.props.loadUserTeams({
          userId: x.userId,
          teamId: x.teamId,
        });
        this.props.messagePush('accepted an invite to a team');
      },
    );
  };

  handleDeclineInvite = () => {
    this.props.patchUserTeamRequest(this.props.myRequest.get('id'), 'no', () =>
      this.props.messagePush('declined an invite to a team')
    );
  };

  handleCancelApplication = () => {
    this.props.patchUserTeamRequest(this.props.myRequest.get('id'), 'no', () =>
      this.props.messagePush('cancelled an application to a team')
    );
  };

  handleApply = () => {
    this.props.createUserTeamRequest({
      userId: this.props.myId,
      teamId: this.props.teamId,
    }, () => this.props.messagePush('applied to a team'));
  };

  handleDisband = () => {
    const { team } = this.props;
    if (!window.confirm(`Disband ${team && team.get('name')}?`)) return;
    this.props.patchTeam(this.props.teamId, 'disband', () =>
      this.props.messagePush('disbanded a team'),
    );
  };

  render() {
    const {
      myId, amAdmin, teamId, team: _team,
      members, myMembership, requests, myRequest, pastMembers, users,
    } = this.props;
    if (!_team) {
      return null;
    }

    if (_team === utils.NOT_FOUND) {
      return <NotFound />;
    }

    const team = _team.toJS();
    if (team.disbandedAt) {
      return (
        <div className={this.cni({ u: 'sectionSingle' })}>
          <div className={this.cn({ u: 'sectionMargined' })}>
            Disbanded team: [{team.abbr}] {team.name}
          </div>
        </div>
      );
    }

    const amLeader = myMembership && myMembership.get('isLeader');
    const amCapable = amLeader || amAdmin;
    const leaderCount = members.reduce((r, x) => // eslint-disable-line
      (x.get('isLeader') ? r + 1 : r), 0,
    );

    let mainButton = null;
    if (myMembership) {
      if (!amLeader || leaderCount > 1) {
        mainButton = (
          <Button
            text="leave" type="important" size="small"
            onClick={this.handleLeave}
          />
        );
      }
    } else if (myRequest) {
      if (myRequest.get('userDecision') === null) {
        mainButton = [
          <Button
            text="accept invite" type="important" size="small"
            onClick={this.handleAcceptInvite}
          />,
          <Button
            text="decline invite" type="important" size="small"
            onClick={this.handleDeclineInvite}
          />,
        ];
      } else {
        mainButton = (
          <Button
            text="cancel application" type="important" size="small"
            onClick={this.handleCancelApplication}
          />
        );
      }
    } else if (myId) {
      mainButton = (
        <Button
          text="apply" type="important" size="small"
          onClick={this.handleApply}
        />
      );
    }
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          Team: [{team.abbr}] {team.name} {mainButton}
          {amAdmin &&
            <Button
              text="disband" type="important" size="small"
              onClick={this.handleDisband}
            />
          }
          {amCapable &&
            <div>
              <Link
                to={`/teams/${teamId}/settings`}
                className={this.cn({ d: 'settings' })}
              >
                SETTINGS
              </Link>
            </div>
          }
          {!!members.size && [
            <div key="members-heading">
              Members
              <span className={this.cn({ d: 'secondary' })}>
                &nbsp;[{members.size}]
              </span>
            </div>,
            <ul key="members">{members.toJS().map(item => (
              <TeamUserTeam
                key={item.id}
                item={item}
                user={users.get(item.userId)}
                isMe={item.userId === myId}
                amLeader={amLeader}
                amAdmin={amAdmin}
                leaderCount={leaderCount}
                patchUserTeam={this.props.patchUserTeam}
                messagePush={this.props.messagePush}
                secondaryClass={this.cn({ d: 'secondary' })}
              />
            ))}</ul>,
          ]}
          {!!requests.size && [
            <div key="requests-heading">
              Requests
              <span className={this.cn({ d: 'secondary' })}>
                &nbsp;[{requests.size}]
              </span>
            </div>,
            <ul key="requests">{requests.toJS().map(item => (
              <TeamUserTeamRequest
                key={item.id}
                item={item}
                user={users.get(item.userId)}
                amLeader={amLeader}
                amAdmin={amAdmin}
                patchUserTeamRequest={this.props.patchUserTeamRequest}
                loadUserTeams={this.props.loadUserTeams}
                messagePush={this.props.messagePush}
                secondaryClass={this.cn({ d: 'secondary' })}
              />
            ))}</ul>,
          ]}
          {!!pastMembers.size && [
            <div key="pastMembers-heading">
              Past members
              <span className={this.cn({ d: 'secondary' })}>
                &nbsp;[{pastMembers.size}]
              </span>
            </div>,
            <ul key="pastMembers">{pastMembers.toJS().map((x) => {
              const user = users.get(x.userId);
              return (
                <li key={x.id}>
                  {user && <Link to={`/users/${user.get('id')}`}>
                    {user.get('nickname')}
                  </Link>}
                  <span className={this.cn({ d: 'secondary' })}>{` /
                    ${moment(x.createdAt).format('YYYY-MM-DD')} ~
                    ${moment(x.leftAt).format('YYYY-MM-DD')} /
                    ${x.kickedBy ? 'kicked' : 'left'}
                  `}</span>
                </li>
              );
            })}</ul>,
          ]}
        </div>
      </div>
    );
  }
}
