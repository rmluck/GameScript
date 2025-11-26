# Standings and Tiebreaker Rules

## NFL Rules

### Regular Season Standings
* Teams are grouped by division (4 teams per division) and conference (4 divisions per conference)
* Win percentage = Wins / (Wins + Losses + 0.5 * Ties)
* Ties count as one-half win and one-half loss for both teams
* Need to calculate and retrieve overall record, win percentage, home record, away record, inter-division record, inter-conference record, points for, points against, points differential, games back, current win/loss streak, strength of schedule, strength of victory, head-to-head records (between tied teams), common opponents (games all tied teams played), playoff seeding, draft order

### Playoff Seeding
* 7 teams per conference make the playoffs
* **Seeds 1-4**: Division winners (ranked by record)
* **Seeds 5-7**: Wild card teams (ranked by record)
* **Seed 1**: Gets first-round bye

### Tiebreaker Procedures

#### To break a tie within a division
* If, at the end of the regular season, two or more teams in the same division finish with identical win percentage, the following steps should be taken until a champion is determined:
    * Two teams:
        1. Head-to-head (best win percentage in games between the two teams)
        2. Best win percentage in games played within the division
        3. Best win percentage in common games
        4. Best win percentage in games played within the conference
        5. Strength of victory in all games
        6. Strength of schedule in all games
        7. Point differential
        8. Points scored
        10. Points allowed
        11. Coin toss
    * Three teams:
        * Note: If two teams remain tied after one or more teams are eliminated during any step, tiebreaker restarts at Step 1 of two-team format. If three teams remain tied after a fourth team is eliminated during any step, tiebreaker restarts at step 1 of three-team format.
        1. Head-to-head (best win percentage in games among the teams)
        2. Best win percentage in games played within the division
        3. Best win percentage in common games
        4. Best win percentage in games played within the conference
        5. Strength of victory in all games
        6. Strength of schedule in all games
        7. Coin toss

#### To break a tie for the wild-card team
* If it is necessary to break ties to determine the three wild card teams from each conference, the following steps will be taken
    1. If the tied teams are from the same division, apply division tiebreaker
    2. If the tied teams are from different divisions, apply the following steps:
    * Two teams:
        1. Head-to-head if applicable
        2. Best win percentage in games played within the conference
        3. Best win percentage in common games, minimum of four
        4. Strength of victory in all games
        5. Strength of schedule in all games
        6. Point differential
        7. Points scored
        8. Points allowed
        9. Coin toss
    * Three or more teams:
        * Note: If two teams remain tied after one or more teams are eliminated during any step, tiebreaker restarts at step 1 of two-team format. If three teams remain tied after fourth team is eliminated during any step, tiebreaker restarts at step 2 of three-team format.
        1. Apply division tiebreaker to eliminate all but highest ranked team in each division prior to proceeding to step 2. Original seeding within division upon application of division tiebreaker remains same for all subsequent applications of procedure that are necessary to identify the two wild card participants.
        2. Head-to-head sweep (applicable only if one team has defeated each of the others or if one team has lost to each of the others)
        3. Best win percentage in games played within conference
        4. Best win percentage in common games, minimum of four
        5. Strength of victory in all games
        6. Strength of schedule in all games
        7. Coin toss
    * When the first wild card team has been identified, the procedure is repeated ot name the second and third wild card (i.e., eliminate all but the highest-ranked team in each division prior to proceeding to step 2). In situations in which three teams from same division are involved in procedure, original seeding of teams remains the same for subsequent applications of the tiebreaker if top-ranked team in that division qualifies for wild card berth.

#### Other tie-breaking procedures
1. Only one team advances to playoffs in any tie-breaking step. Remaining tied teams revert to first step of applicable division or wild card tiebreakers. As an example, if two teams remain tied in any tiebreakaer step after all other teams have been eliminated, the procedure reverts to step 1 of the two-team format to determine the winner. When one team wins the tiebreaker, all other teams revert to step 1 of the applicable two-team or three-team format.
2. In comparing records against common opponents among tied teams, the best win percentage is the deciding factor, since teams may have played an unequal number of games
3. To determine home-field priority among division winners, apply wild card tiebreakers.
4. To determine home-field priority for wild card qualifiers, apply division tiebreakers (if teams are from the same divisiion) or wild card tiebreakers (if teams are from different divisions).
5. To determine the best combined ranking among conference teams in points scored and points allowed, add a team's position in the two categories and the lowest score wins. For example, if Team A is first in points scored and second in points allowed, its combined ranking is 3. If Team B is third in points scored and first in points allowed, its combined ranking is 4. Team A then wins the tiebreaker. If two teams are tied for a position, both teams are awarded the ranking as if they held it solely. For example, if Team A and Team B are tied for first in points scored, each team is assigned a ranking of 1 in that category, and if Team C is third, its ranking will still be 3.

#### Tiebreaker procedure for draft selection meeting
1. Teams not participating in the playoffs shall select in the first through 19th positions in league-wide reverse-standings order.
2. Teams participating in the playoffs shall select according to the following procedures:
    a. The losers of the wild card games shall select in the 19th through 24th positions based on win percentage in reverse-standings order.
    b. The losers of the divisional playoff games shall select in the 25th through 28th positions based on win percentage in reverse-standings order.
    c. The losers of conference championship games shall select 29th and 30th based on win percentage in reverse-standings order.
    d. The winner of the Super Bowl shall select last and the Super Bowl loser will select next-to-last.
3. If ties exist in any grouping, such ties shall be broken by figuring the aggregate win percentage of each involved team's regular season opponents and awarding preferential selection order to the team that faced the schedule of teams with the lowest aggregate win percentage.
4. If ties still exist, apply the divisional, conference, or inter-conference tiebreaking methods, whichever is applicable.
    a. For divisional or conference ties, use procedures on the previous page.
    b. For inter-conference ties, use the following procedures:
        i. Ties involving two teams from different conferences will be broken by (a) head-to-head meeting; (b) best win percentage in common games, minimum of four; (c) strength of victory in all games; (d) best combined ranking among all teams in points scored and points allowed in all games; (e) best net points in all games; (f) best net touchdowns in all games, and finally (g) coin toss.
        ii. Ties involving three or more teams from different conferences will be broken by applying (a) divisional tiebreakers to determine lowest-ranked team in a division, (b) conference tiebreakers to determine lowest ranked team within conference, and (c) interconference tiebreakers to determine the lowest ranked team in the league. The process will be repeated until draft order has been established.

### Other Notes

* **Pct Calculation System**: Treat a 0-0 score calculation as less than 0.000 - so a team that is 0-0 will be lower than a team that is 0-1. This is extremely important in the early season when not many games have been played.
* **Common Games Handling**: The game must have an outcome in order for it to count towards creating a common opponent. At the end of the season, this won't matter but it may impact ranking during the season.
* **Strength of Schedule Methodology**: The game must have an outcome in order for it to count towards a team's Strength of Schedule. At the end of the season, this won't matter but it may impact standings or draft order during the season.

### Algorithm

Pseudocode for standings calculation:

function calculateNFLStandings(scenario):
1. Get all games for season
2. Merge actual results of games with user picks (user picks take priority)
3. Calculate necessary statistics for each team
4. Group teams by division and conference
5. Rank division winners per conference
6. Rank wild cards
7. Apply tiebreakers
8. Return seeded teams + draft order

---

## NBA Rules

### Regular Season Standings

### Playoff Seeding

### Tiebreaker Procedures

### Play-In Tournament

### Draft Lottery & Draft Order

---

## CFB Rules

### Regular Season Rankings

### College Football Playoff

### Conference Standings

### Conference Tiebreakers

### Bowl Game Selection

---

