import type { Game } from '$types';

export interface NBAWeekDateRange {
    week: number;
    startDate: Date;
    endDate: Date;
}

export function getNBAWeekDateRangesFromGames(allGames: Game[]): Map<number, NBAWeekDateRange> {
    const weekRanges = new Map<number, NBAWeekDateRange>();
    
    // Group games by week
    const gamesByWeek = new Map<number, Date[]>();
    allGames.forEach(game => {
        if (game.week) {
            if (!gamesByWeek.has(game.week)) {
                gamesByWeek.set(game.week, []);
            }
            gamesByWeek.get(game.week)!.push(new Date(game.start_time));
        }
    });
    
    // Calculate date ranges for each week
    gamesByWeek.forEach((dates, week) => {
        if (dates.length === 0) return;
        
        // Find the earliest game date for this week
        const earliestGame = new Date(Math.min(...dates.map(d => d.getTime())));
        
        // Find the Monday before (or on) the earliest game
        const startDate = new Date(earliestGame);
        const dayOfWeek = startDate.getDay();
        const daysToSubtract = (dayOfWeek + 6) % 7; // Days back to Monday
        startDate.setDate(startDate.getDate() - daysToSubtract);
        startDate.setHours(0, 0, 0, 0);
        
        // End date is the following Sunday
        const endDate = new Date(startDate);
        endDate.setDate(endDate.getDate() + 6); // Monday + 6 days = Sunday
        endDate.setHours(23, 59, 59, 999);
        
        weekRanges.set(week, {
            week,
            startDate,
            endDate
        });
    });
    
    return weekRanges;
}

export function formatNBADateRange(startDate: Date, endDate: Date): string {
    const options: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' };
    const start = startDate.toLocaleDateString('en-US', options);
    const end = endDate.toLocaleDateString('en-US', options);
    return `${start} - ${end}`;
}

export function getCurrentNBAWeekFromGames(allGames: Game[]): number {
    const now = new Date();
    now.setHours(0, 0, 0, 0);
    
    const weekRanges = getNBAWeekDateRangesFromGames(allGames);
    
    // Check each week to find where current date falls
    for (let week = 1; week <= 25; week++) {
        const range = weekRanges.get(week);
        if (!range) continue;
        
        // Check if we're currently in this week's range
        // Week starts on Monday and ends on Sunday night
        if (now >= range.startDate && now <= range.endDate) {
            return week;
        }
    }
    
    // If we're past all weeks, check if we're between weeks
    // Return the next upcoming week that hasn't started yet
    for (let week = 1; week <= 25; week++) {
        const range = weekRanges.get(week);
        if (!range) continue;
        
        if (now < range.startDate) {
            // We're before this week starts, so return the previous week
            // (or week 1 if this is week 1)
            return Math.max(1, week - 1);
        }
    }
    
    // We're after all regular season weeks
    return 18;
}