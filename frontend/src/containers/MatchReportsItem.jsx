import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

import Button from '../components/Button';

@connect(createStructuredSelector({
  amAdmin: selectors.amAdmin,
  users: selectors.users,
  teams: selectors.teams,
}), actions)
export default class MatchReportsItem extends Component {
  static propTypes = {
    matchReport: ImmutablePropTypes.map,
    publishedReportId: PropTypes.string,
    myTeamId: PropTypes.string,
    matchId: PropTypes.string,
    linkPrefix: PropTypes.string,

    amAdmin: PropTypes.bool,
    users: ImmutablePropTypes.map,
    teams: ImmutablePropTypes.map,

    patchMatchReport: PropTypes.func.isRequired,
    loadMatch: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleAgree = () => {
    const { matchReport } = this.props;
    if (!matchReport) return;
    this.props.patchMatchReport(matchReport.get('id'), 'agree', () => {
      this.props.loadMatch(this.props.matchId);
      this.props.linkPrefix && this.context.router.push(
        this.props.linkPrefix,
      );
      this.props.messagePush('agreed on a match report');
    });
  };

  render() {
    const {
      matchReport, publishedReportId, myTeamId, amAdmin, users, teams,
    } = this.props;
    if (!matchReport) return null;
    const x = matchReport.toJS();
    const teamName = teams.getIn([x.teamBy, 'name']);
    return (
      <li>
        {x.id === publishedReportId && '[PUBLISHED] '}
        [{users.getIn([x.createdBy, 'nickname'])}
        {teamName && ` of ${teamName}`}]
        &nbsp;{x.scoreX} - {x.scoreY} / {x.rawScoreX} - {x.rawScoreY}
        {x.teamBy && ((myTeamId && myTeamId !== x.teamBy) || amAdmin) &&
          <Button
            text="agree" type="important" size="small"
            onClick={this.handleAgree}
          />
        }
      </li>
    );
  }
}
