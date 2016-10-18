import moment from 'moment';
import { PropTypes } from 'react';

import Component from '../utils/Component';

import BracketTeam from './BracketTeam';

const roundWidth = 320;
const matchWidth = 200;
const roundPadding = (roundWidth - matchWidth) / 2;

const teamHeight = 32;
const teamDistance = 2;

const matchHeight = (teamHeight * 2) + teamDistance;

const connectionSegment = 30;
const connectionWidth = roundPadding * 2;

const accWidth = 3;

const dateHeight = 24;

export default class BracketMatch extends Component {
  static propTypes = {
    bracketPath: PropTypes.string.isRequired,

    match: PropTypes.object, // eslint-disable-line
    titleText: PropTypes.string,
    pad: PropTypes.bool,
    hlTeamId: PropTypes.string,
    onTeamMouseEnter: PropTypes.func,
    onTeamMouseLeave: PropTypes.func,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleClick = (e) => {
    if (
      e.nativeEvent.button !== undefined && // not a tap
      e.nativeEvent.button !== 0 // not left click
    ) return;
    this.context.router.push(
      `${this.props.bracketPath}/matches/${this.props.match.id}`,
    );
  };

  renderConnection(isY = false) {
    const { match } = this.props;
    let parent;
    if (!isY) {
      if (match.parentMatchXIsLoser) return null;
      parent = match.parentMatchX;
    } else {
      if (match.parentMatchYIsLoser) return null;
      parent = match.parentMatchY;
    }
    if (!parent) return null;

    let parentY = parent.offset;
    if (parent.winner === 'x') parentY += teamHeight / 2;
    else if (parent.winner === 'y') parentY += matchHeight - (teamHeight / 2);
    else parentY += matchHeight / 2;

    let matchY = match.offset;
    if (!isY) matchY += teamHeight / 2;
    else matchY += matchHeight - (teamHeight / 2);

    let mod = 'pending';
    let stroke = null;
    if (parent.winner) {
      mod = null;
      stroke = 'url(#connectionWin)';
    } else if (parent.teamX && parent.teamY) {
      mod = null;
      stroke = 'url(#connectionSeeded)';
    }

    return (
      <g transform={`translate(-${connectionWidth}, 0)`}>
        <polyline
          className={this.cn({ d: 'connection', m: mod })}
          stroke={stroke}
          points={`
            ${0},${parentY}
            ${connectionSegment},${parentY}
            ${connectionWidth - connectionSegment},${matchY}
            ${connectionWidth},${matchY}
          `}
        />
      </g>
    );
  }

  render() {
    const {
      match, titleText, pad, hlTeamId, onTeamMouseEnter, onTeamMouseLeave,
    } = this.props;
    let mod = 'pending';
    if (match.winner) {
      if (match.winner === 'draw') {
        mod = 'draw';
      } else {
        mod = 'win';
      }
    } else if (match.teamX && match.teamY) {
      mod = 'seeded';
    }

    const isFullySeeded = !!match.teamX && !!match.teamY;

    return (
      <g
        className={this.cni()}
        transform={`translate(
          ${((match.bracketRound - 1) * roundWidth) + roundPadding}, 0
        )`}
      >
        {this.renderConnection()}
        {this.renderConnection(true)}
        <g
          transform={`translate(
            ${pad ? roundPadding / 2 : 0},
            ${match.offset}
          )`}
          onClick={this.handleClick}
          onTouchTap={this.handleClick}
        >
          <rect
            className={this.cn({ d: 'hover' })}
            x={-((accWidth / 2) + teamDistance)}
            y={-((accWidth / 2) + teamDistance)}
            width={matchWidth + accWidth + (teamDistance * 2)}
            height={matchHeight + accWidth + (teamDistance * 2)}
            strokeWidth={accWidth}
          />
          <rect
            className={this.cn({ d: 'date' })}
            x={-(accWidth + teamDistance)}
            y={matchHeight + teamDistance}
            width={matchWidth + ((accWidth + teamDistance) * 2)}
            height={dateHeight}
          />
          <text
            className={this.cn({ d: 'dateText' })}
            x={12 - (accWidth + teamDistance)}
            y={matchHeight + teamDistance + (dateHeight / 2)}
            dy="4" // dominant-baseline not supported by IE
          >
            {moment(match.startedAt).format('MMMM D [at] HH:mm')}
          </text>
          <BracketTeam
            className={this.cn({ d: 'team' })}
            seed={match.seedX}
            teamId={match.teamX}
            score={match.scoreX}
            winner={match.winner}
            isFullySeeded={isFullySeeded}
            titleText={titleText || null}
            hlTeamId={hlTeamId}
            onMouseEnter={onTeamMouseEnter}
            onMouseLeave={onTeamMouseLeave}
          />
          <BracketTeam
            className={this.cn({ d: 'team' })}
            seed={match.seedY}
            teamId={match.teamY}
            score={match.scoreY}
            winner={match.winner}
            isFullySeeded={isFullySeeded}
            isY
            hlTeamId={hlTeamId}
            onMouseEnter={onTeamMouseEnter}
            onMouseLeave={onTeamMouseLeave}
          />
          <rect
            className={this.cn({ d: 'acc', m: mod })}
            x={matchWidth - accWidth}
            y={match.winner !== 'y' ? 0 : teamHeight + teamDistance}
            width={accWidth}
            height={
              (!match.winner || match.winner === 'draw') ?
                matchHeight : teamHeight
            }
          />
        </g>
      </g>
    );
  }
}
