import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as selectors from '../selectors';
import Component from '../utils/Component';

const teamHeight = 32;
const teamDistance = 2;
const accWidth = 3;
const matchWidth = 200;

@connect(createStructuredSelector({
  teams: selectors.teams,
}))
export default class BracketTeam extends Component {
  static propTypes = {
    seed: PropTypes.number,
    teamId: PropTypes.string,
    score: PropTypes.number,
    winner: PropTypes.string,
    isFullySeeded: PropTypes.bool,
    titleText: PropTypes.string,
    isY: PropTypes.bool,
    hlTeamId: PropTypes.string,
    onMouseEnter: PropTypes.func,
    onMouseLeave: PropTypes.func,

    teams: ImmutablePropTypes.map,
  };

  static defaultProps = {
    isFullySeeded: false,
    isY: false,
  };

  handleMouseEnter = () =>
    this.props.onMouseEnter && this.props.onMouseEnter(this.props.teamId);
  handleMouseLeave = () =>
    this.props.onMouseLeave && this.props.onMouseLeave();

  render() {
    const {
      seed, teamId, score, winner, isFullySeeded, titleText, isY, hlTeamId,
      teams,
    } = this.props;

    let name = teamId;
    const team = teams.get(teamId);
    if (team) name = team.get('name');

    let mod = 'pending';
    let fill = null;
    let stroke = 0;
    if (winner) {
      if (winner === 'draw') {
        mod = 'draw';
        fill = 'url(#teamDraw)';
      } else if ((winner === 'x' && !isY) || (winner === 'y' && isY)) {
        mod = 'win';
        fill = 'url(#teamWin)';
      } else {
        mod = 'loss';
        fill = 'url(#teamLoss)';
      }
    } else if (name) {
      if (isFullySeeded) {
        mod = 'fullySeeded';
        fill = 'url(#teamFullySeeded)';
      } else {
        mod = 'seeded';
      }
    } else {
      stroke = accWidth;
    }

    return (
      <g
        className={this.cni({ m: (teamId === hlTeamId) && 'highlight' })}
        transform={`translate(
          0, ${(isY ? teamHeight + teamDistance : 0) + (stroke / 2)}
        )`}
      >
        {titleText &&
          <text
            className={this.cn({ d: 'title' })}
            x={12}
            y={-8}
          >
            {titleText}
          </text>
        }
        {seed &&
          <text
            className={this.cn({ d: 'seed' })}
            x={-12}
            y={teamHeight / 2}
            dy="0.3em"
          >
            {seed}
          </text>
        }
        <g
          transform={`translate(${stroke / 2}, 0)`}
          onMouseEnter={this.handleMouseEnter}
          onMouseLeave={this.handleMouseLeave}
        >
          <rect
            className={this.cn({ d: 'rect', m: mod })}
            fill={fill}
            width={matchWidth - stroke}
            height={teamHeight - stroke}
            strokeWidth={stroke}
          />
          <rect
            className={this.cn({ d: 'hover', m: mod })}
            width={matchWidth - stroke}
            height={teamHeight - stroke}
          />
          <text
            className={this.cn({ d: 'name', m: mod })}
            x="12"
            y={teamHeight / 2}
            dy="0.3em"
            clipPath="url(#teamClip)"
            mask={`url(#teamFade${score === null ? '' : 'More'})`}
          >
            {name}
          </text>
          {score !== null &&
            <text
              className={this.cn({ d: 'score', m: mod })}
              x={matchWidth - 12}
              y={teamHeight / 2}
              dy="0.3em"
            >
              {score}
            </text>
          }
        </g>
      </g>
    );
  }
}
