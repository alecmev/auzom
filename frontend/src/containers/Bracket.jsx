import _ from 'lodash';
import { PropTypes } from 'react';
import ElementPan from 'react-element-pan';
import ImmutablePropTypes from 'react-immutable-proptypes';

import Component from '../utils/Component';

import BracketMatch from './BracketMatch';
import BracketStanding from './BracketStanding';
import BracketTeam from './BracketTeam';

// TODO: explain constants better + reuse stuff with BracketMatch

const bracketPadding = 64;

const roundWidth = 320;
const matchWidth = 200;
const roundPadding = (roundWidth - matchWidth) / 2;

const teamHeight = 32;
const teamDistance = 2;

const matchHeight = (teamHeight * 2) + teamDistance;

const pairInnerDistance = 48;
const pairOuterDistance = 32;

const roundHeaderHeight = 64 + bracketPadding;

const teamFadeLess = 8;
const teamFadeMore = 24;

export default class Bracket extends Component {
  static propTypes = {
    seasonPath: PropTypes.string,
    bracket: ImmutablePropTypes.map,
    bracket2Rounds: ImmutablePropTypes.list,
    bracketMatches: ImmutablePropTypes.list,
    bracket2Standings: ImmutablePropTypes.list,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  state = {
    hlTeamId: null,
  };

  handleTeamMouseEnter = teamId => this.setState({ hlTeamId: teamId });
  handleTeamMouseLeave = () => this.setState({ hlTeamId: null });

  renderGroupStage(rounds, matches) {
    const roundMap = {};
    for (const round of rounds) {
      roundMap[round.number] = round;
    }

    const mpr = {}; // matches per round
    let maxMatches = 0;
    let isMaxOdd = false;
    for (const match of matches) {
      mpr[match.bracketRound] = (mpr[match.bracketRound] + 1) || 1;
      if (mpr[match.bracketRound] > maxMatches) {
        maxMatches = mpr[match.bracketRound];
        isMaxOdd = roundMap[match.bracketRound].byeTeamId !== null;
      } else if (
        mpr[match.bracketRound] === maxMatches &&
        roundMap[match.bracketRound].byeTeamId !== null
      ) {
        isMaxOdd = true;
      }
    }

    const height = (
      (maxMatches * matchHeight) +
      ((maxMatches - 1) * pairOuterDistance) +
      (isMaxOdd ? (pairInnerDistance + teamHeight) : 0)
    );

    const byes = [];
    for (const round of rounds) {
      if (round.byeTeamId !== null) {
        byes.push(
          <g
            transform={`translate(
              ${((round.number - 1) * roundWidth) + roundPadding}, ${
                (mpr[round.number] * matchHeight) +
                ((mpr[round.number] - 1) * pairOuterDistance) +
                pairInnerDistance
              }
            )`}
            key={`bye${round.id}`}
          >
            <BracketTeam
              teamId={round.byeTeamId}
              isFullySeeded
              titleText={'bye'}
              hlTeamId={this.state.hlTeamId}
              onMouseEnter={this.handleTeamMouseEnter}
              onMouseLeave={this.handleTeamMouseLeave}
            />
          </g>
        );
      }
    }

    return {
      height,
      svg: matches.map((x) => {
        x.offset = x.sortNumber * (matchHeight + pairOuterDistance);
        return (
          <BracketMatch
            key={x.id}
            match={x}
            hlTeamId={this.state.hlTeamId}
            onTeamMouseEnter={this.handleTeamMouseEnter}
            onTeamMouseLeave={this.handleTeamMouseLeave}
            {...this.props}
          />
        );
      }).concat(byes),
    };
  }

  renderSingleElimination(rounds, matches) {
    const matchMap = _.keyBy(matches, 'id');
    for (const match of matches) {
      if (match.parentMatchX) {
        match.parentMatchX = matchMap[match.parentMatchX];
      }
      if (match.parentMatchY) {
        match.parentMatchY = matchMap[match.parentMatchY];
      }
    }

    const lastRound = _.maxBy(rounds, 'number').number;
    let goldMatch;
    let bronzeMatch;
    for (const match of matches.filter(x => x.bracketRound === lastRound)) {
      if (!match.parentMatchXIsLoser) goldMatch = match;
      else bronzeMatch = match;
    }

    const roundPairHeight = { 1: (matchHeight * 2) + pairInnerDistance };
    for (let i = 2, j = 1; i <= rounds.length; i += 1, j *= 2) {
      roundPairHeight[i] = (
        (roundPairHeight[1] * j) +
        (pairOuterDistance * j) +
        (teamHeight * 2) +
        teamDistance // two halves of
      );
    }

    let height = bronzeMatch ?
      ((matchHeight * 2) + pairInnerDistance) :
      matchHeight;

    const spaceOut = (match, offset = 0) => {
      let moveParent = 0;
      if (offset < 0) {
        match.offset = 0;
        moveParent = -offset;
        height += -offset;
      } else {
        match.offset = offset;
        if (offset + matchHeight > height) {
          height = offset + matchHeight;
        }
      }

      if (match.parentMatchX) {
        const moveMore = spaceOut(match.parentMatchX, (
            match.offset +
            teamHeight +
            (teamDistance / 2)
          ) - (roundPairHeight[match.parentMatchX.bracketRound] / 2),
        );
        match.offset += moveMore;
        moveParent += moveMore;
      }

      if (match.parentMatchY) {
        spaceOut(match.parentMatchY, (
            match.offset +
            teamHeight +
            (teamDistance / 2) +
            (roundPairHeight[match.parentMatchY.bracketRound] / 2)
          ) - matchHeight,
        );
      }

      return moveParent;
    };

    spaceOut(goldMatch);
    if (bronzeMatch) {
      bronzeMatch.offset = goldMatch.offset + matchHeight + pairInnerDistance;
    }

    return {
      height,
      svg: matches.map((x) => {
        const isBronze = x === bronzeMatch;
        return (
          <BracketMatch
            key={x.id}
            match={x}
            titleText={isBronze ? 'bronze match' : null}
            pad={isBronze}
            hlTeamId={this.state.hlTeamId}
            onTeamMouseEnter={this.handleTeamMouseEnter}
            onTeamMouseLeave={this.handleTeamMouseLeave}
            {...this.props}
          />
        );
      }),
    };
  }

  render() {
    const {
      seasonPath,
      bracket: _bracket,
      bracket2Rounds: _rounds,
      bracketMatches: _matches,
      bracket2Standings: _standings,
    } = this.props;

    if (
      !_bracket ||
      !_rounds || !_rounds.size ||
      !_matches || !_matches.size
    ) return null;

    const bracket = _bracket.toJS();
    const rounds = _rounds.toJS();
    const matches = _matches.toJS();

    for (const match of matches) {
      if (match.scoreX > match.scoreY) match.winner = 'x';
      else if (match.scoreX < match.scoreY) match.winner = 'y';
      else if (match.scoreX !== null) match.winner = 'draw';
      else match.winner = null;
    }

    let body;
    if (
      bracket.type === 'bcl-s8-group-stage' ||
      bracket.type === 'bcl-sc16-swiss' ||
      bracket.type === 'ace-pre-swiss'
    ) {
      body = this.renderGroupStage(rounds, matches);
    } else if (bracket.type === 'bcl-s8-playoffs') {
      body = this.renderSingleElimination(rounds, matches);
    } else {
      return null;
    }

    const width = (roundWidth * rounds.length) + (roundPadding * 2);
    const height = roundHeaderHeight + body.height + bracketPadding;
    const zoom = 1;

    const isPlayoffs = bracket.type === 'bcl-s8-playoffs';
    const isAceSwiss = bracket.type === 'ace-pre-swiss';
    const isSwiss = (
      bracket.type === 'bcl-sc16-swiss' ||
      bracket.type === 'ace-pre-swiss'
    );
    let place = 1;
    let equalsLeft = 0;
    let odd = true;
    const standings = _standings && _standings.map((x) => {
      const equalsBelow = x.get('equalsBelow');
      const row = (
        <BracketStanding
          key={x.get('teamId')}
          standing={x}
          odd={odd}
          hlTeamId={this.state.hlTeamId}
          place={!equalsLeft ? place : null}
          placeTo={equalsBelow ? place + equalsBelow : null}
          playoffs={isPlayoffs}
          swiss={isSwiss}
          aceSwiss={isAceSwiss}
          seasonPath={seasonPath}
          onMouseEnter={this.handleTeamMouseEnter}
          onMouseLeave={this.handleTeamMouseLeave}
        />
      );
      if (equalsBelow) {
        equalsLeft = equalsBelow;
      } else if (equalsLeft) {
        equalsLeft -= 1;
      }

      place += 1;
      odd = !odd;
      return row;
    }).toJS();

    return (
      <div>
        <ElementPan className={this.cn({ d: 'bracket' })}>
          <svg
            className={this.cn({ d: 'bracketInner' })}
            viewBox={`${-roundPadding} 0 ${width} ${height}`}
            width={width * zoom}
            height={height * zoom}
          >
            <defs>
              <linearGradient id="teamFullySeeded">
                <stop
                  className={this.cn({ d: 'gTeamFullySeededA' })}
                  offset="0%"
                />
                <stop
                  className={this.cn({ d: 'gTeamFullySeededB' })}
                  offset="100%"
                />
              </linearGradient>
              <linearGradient id="teamWin">
                <stop
                  className={this.cn({ d: 'gTeamWinA' })}
                  offset="0%"
                />
                <stop
                  className={this.cn({ d: 'gTeamWinB' })}
                  offset="100%"
                />
              </linearGradient>
              <linearGradient id="teamLoss">
                <stop
                  className={this.cn({ d: 'gTeamLossA' })}
                  offset="50%"
                />
                <stop
                  className={this.cn({ d: 'gTeamLossB' })}
                  offset="100%"
                />
              </linearGradient>
              <linearGradient id="teamDraw">
                <stop
                  className={this.cn({ d: 'gTeamDrawA' })}
                  offset="0%"
                />
                <stop
                  className={this.cn({ d: 'gTeamDrawB' })}
                  offset="100%"
                />
              </linearGradient>
              <linearGradient id="connectionSeeded">
                <stop
                  className={this.cn({ d: 'gConnectionSeededA' })}
                  offset="25%"
                />
                <stop
                  className={this.cn({ d: 'gConnectionSeededB' })}
                  offset="75%"
                />
              </linearGradient>
              <linearGradient id="connectionWin">
                <stop
                  className={this.cn({ d: 'gConnectionWinA' })}
                  offset="25%"
                />
                <stop
                  className={this.cn({ d: 'gConnectionWinB' })}
                  offset="75%"
                />
              </linearGradient>
              <linearGradient id="teamFadeGradient">
                <stop
                  stopColor="white"
                  stopOpacity="1"
                  offset="75%"
                />
                <stop
                  stopColor="white"
                  stopOpacity="0"
                  offset="100%"
                />
              </linearGradient>
              <mask id="teamFade">
                <rect
                  width={matchWidth - teamFadeLess}
                  height={teamHeight}
                  fill="url(#teamFadeGradient)"
                />
              </mask>
              <mask id="teamFadeMore">
                <rect
                  width={matchWidth - teamFadeMore}
                  height={teamHeight}
                  fill="url(#teamFadeGradient)"
                />
              </mask>
              <clipPath id="teamClip">
                <rect
                  width={matchWidth}
                  height={teamHeight}
                />
              </clipPath>
            </defs>
            {rounds.map(x => (
              <text
                className={this.cn({ d: 'roundHeader' })}
                x={(x.number - 0.5) * roundWidth}
                y={bracketPadding}
                dy="1em"
                key={x.id}
              >
                {x.name}
              </text>
            ))}
            <g transform={`translate(0, ${roundHeaderHeight})`} key="matches">
              {body.svg}
            </g>
          </svg>
        </ElementPan>
        <div className={this.cn({ d: 'standings', u: 'sectionSingle' })}>
          <div className={this.cn({ u: 'sectionMargined' })}>
            <table>
              <thead>
                <tr>
                  <th className={this.cn({ d: 'standingsPlace' })}>place</th>
                  <th className={this.cn({ d: 'standingsDivider' })} />
                  <th className={this.cn({ d: 'standingsTeam' })}>team</th>
                  {isSwiss && <th>byes</th>}
                  {!isPlayoffs && <th>score</th>}
                  {isAceSwiss && <th title="median-buchholz">mb</th>}
                  {!isPlayoffs && <th title="maps won">mw</th>}
                  {!isPlayoffs && <th title="rounds won">rw</th>}
                  {!isPlayoffs && <th title="raw score ratio">rsr</th>}
                  {!isPlayoffs && <th title="raw score won">rsw</th>}
                  {!isPlayoffs && <th title="raw score lost">rsl</th>}
                  {isPlayoffs && <th title="match loss weight">mlw</th>}
                  {isPlayoffs && <th title="matches lost">ml</th>}
                </tr>
              </thead>
              <tbody>{standings}</tbody>
              <tfoot>
                <tr>
                  <td />
                  <td className={this.cn({ d: 'standingsDividerFoot' })} />
                  <td />
                </tr>
              </tfoot>
            </table>
          </div>
        </div>
      </div>
    );

    // TODO: generate gradients from a config object
  }
}
