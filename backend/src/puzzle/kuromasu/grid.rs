use std::fmt::Write;

use crate::puzzle::kuromasu::{BLOCKED, Cell, FILLED, KuromasuError, MAX_SIZE, UNKNOWN};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub(super) enum Axis {
    Row,
    Col,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub(super) struct Grid {
    size: usize,
    cells: Vec<Cell>,
}

impl Grid {
    pub(super) fn empty(size: usize) -> Self {
        let max_cell = (size * 2 - 2) as Cell;

        Self {
            size,
            cells: vec![max_cell; size * size],
        }
    }

    pub(super) fn new(size: usize, cells: Vec<Cell>) -> Self {
        Self { size, cells }
    }

    pub(super) fn validate_size(size: usize) -> Result<(), KuromasuError> {
        if size == 0 {
            return Err(KuromasuError::InvalidSize);
        }

        if size > MAX_SIZE {
            return Err(KuromasuError::SizeTooLarge(MAX_SIZE));
        }

        Ok(())
    }

    pub(super) fn size(&self) -> usize {
        self.size
    }

    pub(super) fn cells(&self) -> &[Cell] {
        &self.cells
    }

    pub(super) fn cells_mut(&mut self) -> &mut [Cell] {
        &mut self.cells
    }

    pub(super) fn into_cells(self) -> Vec<Cell> {
        self.cells
    }

    pub(super) fn at(&self, row: usize, col: usize) -> Cell {
        self.cells[self.index(row, col)]
    }

    pub(super) fn set(&mut self, row: usize, col: usize, value: Cell) {
        let index = self.index(row, col);
        self.cells[index] = value;
    }

    pub(super) fn sees(&self, row: usize, col: usize) -> usize {
        if self.at(row, col) == BLOCKED {
            return 0;
        }

        self.line_sees(Axis::Row, row, col) + self.line_sees(Axis::Col, col, row)
    }

    pub(super) fn solved(&self) -> bool {
        !self.cells.contains(&UNKNOWN)
    }

    pub(super) fn add_boundary(&mut self, row: usize, col: usize) -> bool {
        self.line_add_boundary(Axis::Row, row, col) | self.line_add_boundary(Axis::Col, col, row)
    }

    pub(super) fn fill(&mut self, row: usize, col: usize) -> usize {
        self.line_fill(Axis::Row, row, col) + self.line_fill(Axis::Col, col, row)
    }

    pub(super) fn open(&self, row: usize, col: usize) -> usize {
        self.line_open(Axis::Row, row, col) + self.line_open(Axis::Col, col, row)
    }

    pub(super) fn valid(&self) -> bool {
        for row in 0..self.size {
            for col in 0..self.size {
                let cell = self.at(row, col);
                if cell >= BLOCKED {
                    continue;
                }

                let sees = self.sees(row, col);
                if sees > usize::from(cell) {
                    return false;
                }

                let open = self.open(row, col);
                if sees + open < usize::from(cell) {
                    return false;
                }
            }
        }

        true
    }

    pub(super) fn render(&self, mask: Option<&[Cell]>) -> Result<String, KuromasuError> {
        if let Some(mask) = mask {
            self.validate_mask(mask)?;
        }

        let mut rendered = String::new();
        write!(rendered, "\n Size: {size} * {size}", size = self.size)
            .expect("writing to a String cannot fail");

        let mut divider = "\n ".to_owned();
        for _ in 0..self.size {
            divider.push_str("-----");
        }
        divider.push_str("-\n");

        for row in 0..self.size {
            rendered.push_str(&divider);
            for col in 0..self.size {
                let index = self.index(row, col);
                let mut value = if self.cells[index] == BLOCKED {
                    "XX".to_owned()
                } else {
                    format!("{:02}", self.cells[index])
                };

                if mask.is_some_and(|mask| mask[index] == 0) {
                    value = "  ".to_owned();
                }

                rendered.push_str(" | ");
                rendered.push_str(&value);
            }
            rendered.push_str(" |");
        }
        rendered.push_str(&divider);

        Ok(rendered)
    }

    pub(super) fn validate_mask(&self, mask: &[Cell]) -> Result<(), KuromasuError> {
        if mask.len() != self.size * self.size {
            return Err(KuromasuError::InvalidMaskLength);
        }

        if mask.iter().any(|value| !matches!(value, 0 | 1)) {
            return Err(KuromasuError::InvalidMaskValue);
        }

        Ok(())
    }

    pub(super) fn masked_board(&self, mask: &[Cell]) -> Result<Vec<Cell>, KuromasuError> {
        self.validate_mask(mask)?;

        Ok(mask
            .iter()
            .zip(&self.cells)
            .map(|(visible, cell)| if *visible == 1 { *cell } else { UNKNOWN })
            .collect())
    }

    pub(super) fn same_shape(&self, solution: &[Cell]) -> bool {
        self.cells.len() == solution.len()
            && self
                .cells
                .iter()
                .zip(solution)
                .all(|(cell, solved)| (*cell == BLOCKED) == (*solved == BLOCKED))
    }

    pub(super) fn line_in_bounds(&self, pos: isize) -> bool {
        pos >= 0 && pos < self.size as isize
    }

    pub(super) fn line_sees_left(&self, axis: Axis, index: usize, pos: usize) -> usize {
        if self.line_at(axis, index, pos) == BLOCKED {
            return 0;
        }

        let mut amount = 0;
        for current in (0..pos).rev() {
            let cell = self.line_at(axis, index, current);
            if cell == BLOCKED || cell == UNKNOWN {
                break;
            }

            amount += 1;
        }

        amount
    }

    pub(super) fn line_sees_right(&self, axis: Axis, index: usize, pos: usize) -> usize {
        if self.line_at(axis, index, pos) == BLOCKED {
            return 0;
        }

        let mut amount = 0;
        for current in pos + 1..self.size {
            let cell = self.line_at(axis, index, current);
            if cell == BLOCKED || cell == UNKNOWN {
                break;
            }

            amount += 1;
        }

        amount
    }

    pub(super) fn line_open_left(&self, axis: Axis, index: usize, pos: usize) -> usize {
        if self.line_at(axis, index, pos) == BLOCKED {
            return 0;
        }

        let mut amount = 0;
        for current in (0..pos).rev() {
            let cell = self.line_at(axis, index, current);
            if cell == BLOCKED {
                break;
            }

            if cell == UNKNOWN {
                amount += 1;
            }
        }

        amount
    }

    pub(super) fn line_open_right(&self, axis: Axis, index: usize, pos: usize) -> usize {
        if self.line_at(axis, index, pos) == BLOCKED {
            return 0;
        }

        let mut amount = 0;
        for current in pos + 1..self.size {
            let cell = self.line_at(axis, index, current);
            if cell == BLOCKED {
                break;
            }

            if cell == UNKNOWN {
                amount += 1;
            }
        }

        amount
    }

    pub(super) fn line_fill_left(
        &mut self,
        axis: Axis,
        index: usize,
        pos: usize,
        max: usize,
    ) -> usize {
        if self.line_at(axis, index, pos) == BLOCKED {
            return 0;
        }

        let mut current = 0;
        let mut filled = 0;
        for line_pos in (0..pos).rev() {
            if max != 0 && current == max {
                break;
            }

            let cell = self.line_at(axis, index, line_pos);
            if cell == BLOCKED {
                break;
            }

            current += 1;

            if cell != UNKNOWN {
                continue;
            }

            self.line_set(axis, index, line_pos, FILLED);
            filled += 1;
        }

        filled
    }

    pub(super) fn line_fill_right(
        &mut self,
        axis: Axis,
        index: usize,
        pos: usize,
        max: usize,
    ) -> usize {
        if self.line_at(axis, index, pos) == BLOCKED {
            return 0;
        }

        let mut current = 0;
        let mut filled = 0;
        for line_pos in pos + 1..self.size {
            if max != 0 && current == max {
                break;
            }

            let cell = self.line_at(axis, index, line_pos);
            if cell == BLOCKED {
                break;
            }

            current += 1;

            if cell != UNKNOWN {
                continue;
            }

            self.line_set(axis, index, line_pos, FILLED);
            filled += 1;
        }

        filled
    }

    pub(super) fn line_at(&self, axis: Axis, index: usize, pos: usize) -> Cell {
        match axis {
            Axis::Row => self.at(index, pos),
            Axis::Col => self.at(pos, index),
        }
    }

    pub(super) fn line_set(&mut self, axis: Axis, index: usize, pos: usize, value: Cell) {
        match axis {
            Axis::Row => self.set(index, pos, value),
            Axis::Col => self.set(pos, index, value),
        }
    }

    fn index(&self, row: usize, col: usize) -> usize {
        row * self.size + col
    }

    fn line_sees(&self, axis: Axis, index: usize, pos: usize) -> usize {
        self.line_sees_left(axis, index, pos) + self.line_sees_right(axis, index, pos)
    }

    fn line_open(&self, axis: Axis, index: usize, pos: usize) -> usize {
        self.line_open_left(axis, index, pos) + self.line_open_right(axis, index, pos)
    }

    fn line_add_boundary(&mut self, axis: Axis, index: usize, pos: usize) -> bool {
        if self.line_at(axis, index, pos) == BLOCKED {
            return false;
        }

        let mut changed = false;

        for line_pos in (0..pos).rev() {
            let cell = self.line_at(axis, index, line_pos);
            if cell == BLOCKED {
                break;
            }

            if cell != UNKNOWN {
                continue;
            }

            self.line_set(axis, index, line_pos, BLOCKED);
            changed = true;
            break;
        }

        for line_pos in pos + 1..self.size {
            let cell = self.line_at(axis, index, line_pos);
            if cell == BLOCKED {
                break;
            }

            if cell != UNKNOWN {
                continue;
            }

            self.line_set(axis, index, line_pos, BLOCKED);
            changed = true;
            break;
        }

        changed
    }

    fn line_fill(&mut self, axis: Axis, index: usize, pos: usize) -> usize {
        self.line_fill_left(axis, index, pos, 0) + self.line_fill_right(axis, index, pos, 0)
    }
}
