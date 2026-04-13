use crate::puzzle::takuzu::TakuzuCell;
use crate::puzzle::takuzu::TakuzuCell::{X, O, Empty};


#[derive(Clone, Copy, Default)]
struct CellCounts {
    x_count: u8,
    o_count: u8,
}

impl CellCounts {
    fn add(&mut self, value: TakuzuCell) {
        match value {
            X => self.x_count += 1,
            O => self.o_count += 1,
            _ => {}
        }
    }

    fn remove(&mut self, value: TakuzuCell) {
        match value {
            X => self.x_count -= 1,
            O => self.o_count -= 1,
            _ => {}
        }
    }

    fn can_balance_with_remaining(&self, half: u8) -> bool {
        self.x_count <= half
            && self.o_count <= half
    }
}

#[derive(Clone, Copy, Default)]
struct LineState {
    counts: CellCounts,
    bits: u32,
}


impl LineState {
    fn add(&mut self, value: TakuzuCell, bit_idx: usize) {
        self.counts.add(value);
        if value == X {
            self.bits |= 1 << bit_idx;
        }
    }

    fn remove(&mut self, value: TakuzuCell, bit_idx: usize) {
        self.counts.remove(value);
        if value == X {
            self.bits &= !(1 << bit_idx);
        }
    }
}

pub(super) fn generate_grid(size: u8) -> Vec<Vec<TakuzuCell>> {
    let mut grid = vec![vec![O; size as usize]; size as usize];

    let mut rows = generate_possible_rows(size);

    // Optional shuffle
    use rand::seq::SliceRandom;
    let mut rng = rand::thread_rng();
    rows.shuffle(&mut rng);

    let mut used_rows = vec![false; rows.len()];
    let mut states = vec![LineState::default(); size as usize];

    fill_rows(
        &mut grid,
        &rows,
        &mut used_rows,
        &mut states,
        0,
        size,
    );

    grid
}

fn fill_rows(
    grid: &mut Vec<Vec<TakuzuCell>>,
    candidates: &[Vec<TakuzuCell>],
    used_rows: &mut [bool],
    states: &mut [LineState],
    row_idx: usize,
    size: u8,
) -> bool {
    if row_idx == size as usize {
        return true;
    }

    for (idx, candidate) in candidates.iter().enumerate() {
        if used_rows[idx] {
            continue;
        }

        if !can_place_row(grid, states, row_idx, candidate, size) {
            continue;
        }

        grid[row_idx] = candidate.clone();
        apply_row(states, candidate, row_idx);
        used_rows[idx] = true;

        let mut ok = true;

        if row_idx == size as usize - 1 && !columns_unique(states) {
            ok = false;
        }

        if ok && fill_rows(grid, candidates, used_rows, states, row_idx + 1, size) {
            return true;
        }

        used_rows[idx] = false;
        undo_row(states, candidate, row_idx);
    }

    false
}

fn can_place_row(
    grid: &[Vec<TakuzuCell>],
    states: &[LineState],
    row_idx: usize,
    candidate: &[TakuzuCell],
    size: u8,
) -> bool {
    let half = size / 2;

    for (col, state) in states.iter().enumerate() {
        let bit = candidate[col];

        let mut counts = state.counts;
        counts.add(bit);

        if !counts.can_balance_with_remaining(half) {
            return false;
        }

        // No vertical triplets
        if row_idx >= 2 {
            let a = grid[row_idx - 2][col];
            let b = grid[row_idx - 1][col];
            if bit == a && bit == b {
                return false;
            }
        }
    }

    true
}

fn apply_row(states: &mut [LineState], candidate: &[TakuzuCell], row_idx: usize) {
    for (col, &cell) in candidate.iter().enumerate() {
        states[col].add(cell, row_idx);
    }
}

fn undo_row(states: &mut [LineState], candidate: &[TakuzuCell], row_idx: usize) {
    for (col, &cell) in candidate.iter().enumerate() {
        states[col].remove(cell, row_idx);
    }
}

fn columns_unique(states: &[LineState]) -> bool {
    use std::collections::HashSet;

    let mut seen = HashSet::with_capacity(states.len());

    for state in states {
        if !seen.insert(state.bits) {
            return false;
        }
    }

    true
}
fn generate_possible_rows(size: u8) -> Vec<Vec<TakuzuCell>> {
    let mut possible_rows: Vec<Vec<TakuzuCell>> = vec![];


    let mut current_row: Vec<TakuzuCell> = vec![Empty; size as usize];
    fn generate_row (
        row_index: usize,
        current_x_count: u8,
        current_o_count: u8,
        current_row: &mut Vec<TakuzuCell>,
        possible_rows: &mut Vec<Vec<TakuzuCell>>,
        size: u8
    ) {
        let half_size = size / 2;
        if current_o_count > half_size || current_x_count > half_size {
            return;
        }

        if row_index >= 3 {
            if current_row[row_index - 3] == current_row[row_index - 2]
                && current_row[row_index - 2] == current_row[row_index - 1] {
                return;
            }
        }

        if row_index == usize::from(size) {
            if current_x_count == half_size && current_o_count == half_size {
                possible_rows.push(current_row.to_vec());
            }
            return;
        }

        current_row[row_index] = X;
        generate_row(row_index + 1, current_x_count + 1, current_o_count, current_row, possible_rows, size);

        current_row[row_index] = O;
        generate_row(row_index + 1, current_x_count, current_o_count + 1, current_row, possible_rows, size);
    }
    generate_row(0, 0, 0, &mut current_row, &mut possible_rows, size);
    possible_rows
}